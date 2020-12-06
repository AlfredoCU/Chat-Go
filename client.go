package main

import (
	"bufio"
	"container/list"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

var name string
var messageClient list.List
var showClient = false

type FileClient struct {
	BS []byte
	Name string
	UserName string
}

func connectClient() {
	c, err := net.Dial("tcp", ":9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	err = gob.NewEncoder(c).Encode(name)

	if err != nil {
		fmt.Println(err)
	}

	go receiveMessageClient(c)

	Menu(c)
	_ = c.Close()
}

func Menu(c net.Conn) {
	var op string
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("MENU")
		fmt.Println("1.- Send Message")
		fmt.Println("2.- Send File")
		fmt.Println("3.- Show Message")
		fmt.Println("0.- Exit")
		fmt.Print("Option: ")
		input.Scan()
		op = input.Text()

		if op == "1" {
			var msg string
			fmt.Print("Message: ")
			input.Scan()
			msg = input.Text()
			var submenu uint64 = 1
			err := gob.NewEncoder(c).Encode(submenu)
			if err != nil {
				fmt.Println(err)
			} else {
				messageClient.PushBack("Tú: "+msg)
				msg := name + ": " + msg
				_ = gob.NewEncoder(c).Encode(msg)
			}
		} else if op == "2" {
			var msg string
			fmt.Print("Ruta: ")
			input.Scan()
			msg = input.Text()
			sendFileClient(c, msg)
		} else if op == "3" {
			showClient = true
			showMessageClient()
			input.Scan()
			showClient = false
		} else if op == "0" {
			var submenu uint64 = 3
			err := gob.NewEncoder(c).Encode(submenu)
			if err != nil {
				fmt.Println(err)
			}
			break
		} else {
			fmt.Println("Option invalid!")
		}
	}
}

func showMessageClient() {
	fmt.Println()
	for e:=messageClient.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func receiveMessageClient(c net.Conn) {
	var op uint64
	var msg string
	for {
		err := gob.NewDecoder(c).Decode(&op)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if op == 1 {
			err = gob.NewDecoder(c).Decode(&msg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			messageClient.PushBack(msg)
			if showClient {
				fmt.Println(msg)
			} else {
				showMessageClient()
			}
		} else if op == 2 {
			receiveFileClient(c)
		}
	}
}

func sendFileClient(c net.Conn, route string) {
	file, err := os.Open(route)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// To see the status of the file (size, etc)
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	total := stat.Size() // content size
	bs := make([]byte, total) // Slice de bytes
	count, err := file.Read(bs)
	if err != nil{
		fmt.Println(err,count)
		return
	}
	// Send slice de bytes
	nameFile := file.Name()

	var submenu uint64 = 2
	err = gob.NewEncoder(c).Encode(submenu)
	if err != nil {
		fmt.Println(err)
		return
	}
	f := FileClient{BS: bs, Name: nameFile, UserName: name}
	stopSendFile(c, f)
}

func stopSendFile(c net.Conn, f FileClient) {
	err := gob.NewEncoder(c).Encode(&f)
	if err != nil {
		fmt.Println(err)
	} else {
		messageClient.PushBack("Tú: "+f.Name)
	}
}

// Here it is redirected when the flag indicates file
func receiveFileClient(c net.Conn) {
	var f FileClient
	err := gob.NewDecoder(c).Decode(&f)
	if err != nil {
		//fmt.Println(err)
		return
	}
	// Save file
	file, err := os.Create(f.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, _ = file.Write(f.BS)

	// Save message
	msg := f.UserName + ": " +f.Name
	messageClient.PushBack(msg)
	if showClient {
		fmt.Println(msg)
	} else {
		showMessageClient()
	}
}

func main() {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Name: ")
	input.Scan()
	name = input.Text()
	connectClient()
}
