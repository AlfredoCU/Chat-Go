package main

// Import libraries.
import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

// User struct.
type User struct {
	Name  string
	Message string
}

// Type and Port.
const (
	TypeConn = "tcp"
	PortConn = ":9999"
)

// Login and message chan.
var login = make(chan string)
var messageClient = make(chan User)

/* Function client.
@param login chan string
@param messageClient chan User
*/
func client(login chan string, messageClient chan User) {
	// Received Message.
	var receivedMessage string

	// Connection to server.
	con, err := net.Dial(TypeConn, PortConn)

	// Error.
	if err != nil {
		fmt.Sprintln(err)
		return
	}

	// Close connection.
	defer con.Close()

	// Login, send and received messages.
	for {
		// Actions.
		select {
			// Login to server
			case loginMessage := <-login:
				err = gob.NewEncoder(con).Encode(loginMessage)
				if err != nil {
					fmt.Sprintln(err)
					return
				}

			// Send message.
			case newMessage := <-messageClient:
				chatMessage := newMessage.Name + ": " + newMessage.Message
				err = gob.NewEncoder(con).Encode(chatMessage)
				if err != nil {
					fmt.Sprintln(err)
					return
				}
		}

		// Decode received messages.
		err := gob.NewDecoder(con).Decode(&receivedMessage)

		// Error.
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print Message.
		fmt.Println(receivedMessage)
	}
}

// Function main.
func main() {
	// Variables opc and newUser.
	var opc string
	var newUser User

	// UserName.
	input := bufio.NewScanner(os.Stdin)

	// Start connection.
	go client(login, messageClient)

	// Print userName.
	fmt.Print("Username: ")
	input.Scan()
	nameUser := input.Text()

	// Added in User struct
	newUser.Name = nameUser

	// Connect userName.
	login <- "Has connected: " + nameUser

	// Menu.
	fmt.Println("----------MENU----------")
	fmt.Println("1.- Send Message")
	fmt.Println("2.- Send File")
	fmt.Println("3.- Exit System")

	for {
		// Added options.
		fmt.Print("Option: ")
		_, _ = fmt.Scanln(&opc)

		// Options.
		switch opc {
			// Send messages.
			case "1":
				input.Scan()
				msg := input.Text()
				newUser.Message = msg
				messageClient <- newUser
			// Send files.
			case "2":
				fmt.Println("Not implemented!")
			// Disconnected.
			case "3":
				login <- "It has disconnected: " + nameUser
				var exit string
				_, _ = fmt.Scanln(&exit)
				return
			// Default.
			default:
				invalidOptionsClient()
		}
	}
}

// Function invalidOptionsClient.
func invalidOptionsClient() {
	fmt.Print("\n-Invalid Option!\n\n")
}