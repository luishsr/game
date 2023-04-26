package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	TYPE            = "tcp"
	LIST_PLAYERS    = "list_players"
	REGISTER_PLAYER = "register"
	NO_PLAYERS      = "noplayers"
	START_GAME      = "startgame"
	SEND_GUESS      = "sendguess"
	ASK_FOR_A_GUESS = "Type your guess: "
)

var host string
var port string

func main() {

	host = promptMessage("Type the IP Address: ")
	port = promptMessage("Type the Port number: ")
	player_name := promptMessage("Type your Name: ")

	log.Println("Connecting to " + host + ":" + port + " ...")

	tcpServer, err := net.ResolveTCPAddr(TYPE, host+":"+port)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	//Register as a player into the server
	sendMessageToServer(conn, REGISTER_PLAYER+player_name)

	//closes connection
	conn.Close()

	log.Println("Client: Connection closed")

	for {
		handleMessages()
	}

}

func sendMessageToServer(conn *net.TCPConn, message string) {
	//Sends a message to the TCP server
	_, err := conn.Write([]byte(message))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	log.Println("Client: message sent - " + message)
}

func connectAndSendMessageToServer(message string) {

	tcpServer, err := net.ResolveTCPAddr(TYPE, host+":"+port)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	//Sends a message to the TCP server
	sendMessageToServer(conn, message)

	log.Println("Client: message sent - " + message)
}

func handleMessages() {

	//Start listening to connections
	listen, err := net.Listen(TYPE, host+":"+port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		//Processes messages from the server
		handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Client: received message:" + string(buffer[:]))

	//route requests according to what was sent by the server
	serverCommand := string(buffer[:])

	//process command
	processCommand(serverCommand, conn)

}

func processCommand(command string, conn net.Conn) {

	//Identify command
	if command == SEND_GUESS {

		//Asks the player for a guess
		guess := promptMessage(ASK_FOR_A_GUESS)

		//Send guess to the server
		connectAndSendMessageToServer(guess)
	}
}

func promptMessage(prompt string) string {
	//ask the user for a guess
	fmt.Print(prompt)

	//creates a reader from the keyboard input
	reader := bufio.NewReader(os.Stdin)

	//reads the input
	input, err := reader.ReadString('\n')

	//check for error when reading
	if err != nil {
		log.Fatal(err)
	}

	//remove spaces
	input = strings.TrimSpace(input)

	//returns user input
	return input
}
