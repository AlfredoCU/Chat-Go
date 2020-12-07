package main

// Import libraries.
import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

// Data.
var (
	clientList []net.Conn
	chat []string
)

// Type and Port.
const (
	Type = "tcp"
	Port = ":9999"
)

// Function server.
func server() {
	// Start server.
	serv, err := net.Listen(Type, Port)

	// Error.
	if err != nil {
		fmt.Sprintln("Error: ", err)
		return
	}

	// Accept client or error.
	for {
		// Accept client.
		con, err := serv.Accept()

		// `Error`.
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		// Added Client in array.
		clientList = append(clientList, con)

		// Request of client.
		go handleClient(con)
	}
}

/* Function handleClient.
@param conn net.Conn
*/
func handleClient(con net.Conn) {
	// Message.
	var message string

	// Send Message.
	for {
		// Decode Message.
		err := gob.NewDecoder(con).Decode(&message)

		// Error.
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		// Print message.
		fmt.Println(message)

		// Added message in chat.
		chat = append(chat, message)

		// Disconnected client.
		logoutFlag := strings.Contains(message, "It has disconnected")

		// Increment or decrement of clients.
		if logoutFlag == true {
			for i := 0; i < len(clientList); i++ {
				if con == clientList[i] {
					clientList = append(clientList[:i], clientList[i+1:]...)
				}
			}
		}

		// New Message.
		for i := 0; i < len(clientList); i++ {
			if con != clientList[i] {
				err := gob.NewEncoder(clientList[i]).Encode(message)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

// Function saveMessage
func saveMessage() {
	// Create file.
	file, err := os.Create("Backup_Message.txt")

	// Error.
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// Close file.
	defer file.Close()

	// Added messages in file.
	for _, message := range chat {
		_, _ = file.WriteString(message + "\n")
	}

	// Successfully.
	fmt.Println("Copy of MESSAGES successfully saved.")
}

// Function main.
func main() {
	// Options.
	var opc string

	// Start Server.
	go server()

	// Basic menu.
	fmt.Println("Start Server...")
	fmt.Println("1.- Save messages. 2.- Exit server.")
	fmt.Println("Messages: ")

	for {
		// Option.
		_, _ = fmt.Scanln(&opc)

		switch opc {
			// Save all chat.
			case "1":
				saveMessage()
			// Exit.
			case "2":
				exited()
				return
			// Default.
			default:
				invalidOptions()
			}
	}
}

// Function exited.
func exited() {
	fmt.Println("\n-System exited...")
}

// Function invalidOptions.
func invalidOptions() {
	fmt.Print("\n-Invalid Option!\n\n")
}