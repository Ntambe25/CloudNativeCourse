// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//type client chan<- string // an outgoing message channel

// A struct of type client
// A channel called Channel is used to store clients messages
// Name is used to store clients name
type client struct {
	Channel chan<- string
	Name    string
}

// entering is a channel of type client that stores information on the client that has entered the chat
// leaving  is a channel of type client that stores information on the client that has left    the chat
// messages store messages by the clients
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {

	//Create a server
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	//Function that handles messages, entering and leaving of clients
	go broadcaster()

	//In a while loop, accepts the clients looking to enter the chat
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		//Function Used to fill in the channels with appropriate
		//text for broadcaster to display it
		go handleConn(conn)
	}
}

// Function that handles messages, entering and leaving of clients
func broadcaster() {

	//Clients is a map that contains all current users
	clients := make(map[client]bool)

	//While loop
	for {
		select {

		//If a message has been sent, broadcast the message
		//to all clients message channels
		case msg := <-messages:

			//Sending messages to channels of all the client
			//present in the map clients
			for cli := range clients {
				cli.Channel <- msg
			}

		//If a new client has entered the chat, enter them
		//into the clients map. Display the new users name
		//in their chat and displays the users connected
		case cli := <-entering:
			clients[cli] = true

			cli.Channel <- "Welcome " + cli.Name + "\nCurrent Clients :"
			// store all connected client's name in current_clients
			for current_clients := range clients {
				cli.Channel <- current_clients.Name
			}

		//If the client has left, remove their entry for the clients map
		//and close their channel
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.Channel)
		}
	}
}

// Function Used to fill in the channels with appropriate
// text for broadcaster to display it
func handleConn(conn net.Conn) {

	//Create a channel ch for outgoing client messages
	ch := make(chan string)

	//Create a empty client struct called temp
	//This is used to
	temp := client{}

	//Function used to print a message to all client channels
	go clientWriter(conn, ch)

	//Get user name
	who := conn.RemoteAddr().String()
	temp.Name = who

	ch <- "You are " + who
	temp.Channel = ch

	messages <- who + " has arrived"
	entering <- temp
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()
	leaving <- temp
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}
