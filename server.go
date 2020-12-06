package main

import (
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

const (
	Type = "tcp"
	Port = ":9999"
)

var (
	message list.List
	client list.List
	show = true
)

type File struct {
	BS []byte
	Name string
	UserName string
}

func server() {
	s, err := net.Listen(Type, Port)

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		c, err := s.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	var op uint64
	var err error

	newClient(c)

	for {
		err = gob.NewDecoder(c).Decode(&op)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if op == 1 {
			receiveMessage(c)
		} else if op == 2 {
			receiveFile(c)
		} else if op == 3 {
			disconnectClient(c)
			return
		} else if op == 0 {
			newClient(c)
		}
	}
}

func newClient(c net.Conn) {
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)

	if err != nil {
		fmt.Println(err)
	}

	client.PushBack(c)
	msg = "Online: " + msg

	if show {
		fmt.Println(msg)
	}

	message.PushBack(msg)
}

func receiveMessage(c net.Conn) {
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	if show {
		fmt.Println(msg)
	}

	message.PushBack(msg)
	sendMessageToEveryone(msg, c)
}

func sendMessageToEveryone(msg string, c net.Conn) {
	var op uint64 = 1

	for e:=client.Front(); e != nil; e = e.Next() {
		if e.Value.(net.Conn) != c {
			err := gob.NewEncoder(e.Value.(net.Conn)).Encode(op)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = gob.NewEncoder(e.Value.(net.Conn)).Encode(msg)

			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func receiveFile(c net.Conn) {
	var f File
	err := gob.NewDecoder(c).Decode(&f)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Save File
	file, err := os.Create(f.Name)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	_, _ = file.Write(f.BS)

	msg := f.UserName+": "+f.Name
	message.PushBack(msg)

	if show {
		fmt.Println(msg)
	}

	// Send file to everyone
	sendFilesAll(c, f)
}

func sendFilesAll(c net.Conn, f File) {
	var op uint64 = 2
	for e:=client.Front(); e != nil; e = e.Next() {
		if e.Value.(net.Conn) != c {
			err := gob.NewEncoder(e.Value.(net.Conn)).Encode(op)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = gob.NewEncoder(e.Value.(net.Conn)).Encode(f)

			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func showMessage() {
	for e:=message.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func disconnectClient(c net.Conn) {
	for e:=client.Front(); e != nil; e = e.Next() {
		if e.Value.(net.Conn) == c {
			client.Remove(e)
			return
		}
	}
}

func backupMessage() {
	// Backup to txt file
	file, err := os.Create("Backup_Message.txt")
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	for e:=message.Front(); e != nil; e = e.Next() {
		_, _ = file.WriteString(e.Value.(string) + "\n")
	}
}

func main() {
	go server()
	var ops int64
	for {
		fmt.Println("MENU")
		fmt.Println("1.- Show Message")
		fmt.Println("2.- Save Message")
		fmt.Println("0.- Exit")
		fmt.Print("Option: ")
		_, _ = fmt.Scanln(&ops)

		if ops == 1 {
			show = true
			showMessage()
			_, _ = fmt.Scanln(&ops)
			show = false
		} else if ops == 2 {
			backupMessage()
		} else if ops == 0 {
			break
		} else {
			fmt.Println("Option invalid!")
		}
	}
}