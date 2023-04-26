package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

// stores the target number for the guess game
var targetNumber int = rand.Intn(100)

// max guesses
var maxGuesses int = 10

// count guesses
var count int

const (
	HOST            = "localhost"
	PORT            = "8081"
	TYPE            = "tcp"
	CONNECT         = "connect"
	LIST_PLAYERS    = "list_players"
	REGISTER_PLAYER = "register"
	NO_PLAYERS      = "noplayers"
	START_GAME      = "startgame"
	SEND_GUESS      = "sendguess"
)

type Player struct {
	ipAddress string
	name      string
	hasPlayed bool
}

// List of registered players
var players []Player

func main() {

	//Start listening to connections
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
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

		//Processes messages from clients
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// incoming request
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	//route requests according to what was sent by the client
	clientCommand := string(buffer[:])

	log.Println("Server: message receive - " + clientCommand)

	//process command and respond to client
	responseStr := processCommand(clientCommand, conn)

	// write data to response
	conn.Write([]byte(responseStr))

}

func processCommand(command string, conn net.Conn) string {

	//Stores response
	var response string = ""

	//Identify command
	if command == REGISTER_PLAYER {

		//Gets player details and client address
		newPlayer := Player{conn.RemoteAddr().String(), substr(command, 8, len(command)), false}

		//Adds the player to the list of players
		players = append(players, newPlayer)

		//Checks if there are other players available
		if len(players) < 2 {
			response = NO_PLAYERS
		} else {
			response = START_GAME

			//Starts the game
			startGame()
		}
	}

	return response
}

// Starts the game
func startGame() {

	//the output message to the user
	var message string

	//repeat the loop until run out of guesses or guess it right
	for count = 1; count <= maxGuesses; count++ {

		//Picks a player
		player := pickThePlayer()

		//Get players guess
		input := messagePlayer(player, SEND_GUESS)

		//remove spaces
		input = strings.TrimSpace(input)

		//converts the read string to int
		guess, err := strconv.Atoi(input)

		//check for errors
		if err != nil {
			//log error
			log.Fatal(err)
		}

		//check if the guess was lower than the target
		if guess < targetNumber && count < maxGuesses {
			message = "Oops, your guess was LOW. You have " + strconv.Itoa(maxGuesses-count) + " guesses left."
		} else if guess > targetNumber && count < maxGuesses {
			message = "Oops, your guess was HIGH. You have " + strconv.Itoa(maxGuesses-count) + " guesses left."
		} else if count < maxGuesses {
			message = "Good job! You guessed it!"

			//Message the player
			messagePlayer(player, message)

			//breaks the loop
			break

		} else {
			//run out of guesses
			message = "Sorry, you didn't guess my number. It was " + strconv.Itoa(targetNumber)

			//Message the player
			messagePlayer(player, message)
		}

		//print the message
		messagePlayer(player, message)
	}

}

// Pick a player
func pickThePlayer() Player {

	//Stores the picked player
	var playerToPlay Player
	var theOtherPlayer Player

	//Select a player according to who played last (or the first one in the list on the first play)
	playerA := players[0]

	//Played already?
	if playerA.hasPlayed {
		//Picks the other
		playerToPlay = players[1]
		//Sets the other play
		theOtherPlayer = players[0]

	} else {
		playerToPlay = playerA
		//Sets the other play
		theOtherPlayer = players[1]
	}

	//Marks player as hasPlayed
	playerToPlay.hasPlayed = true

	//Mark the other player as not played
	theOtherPlayer.hasPlayed = false

	return playerToPlay
}

// Sends a message to a Player and collect response
func messagePlayer(player Player, message string) string {

	//Establish a connection with the player
	tcpServer, err := net.ResolveTCPAddr(TYPE, player.ipAddress)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	//Send the message and returns response
	response := sendAndReceive(conn, SEND_GUESS)

	return response
}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func sendAndReceive(conn *net.TCPConn, message string) string {
	//Sends a message to the TCP server
	_, err := conn.Write([]byte(message))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	log.Println("Server: message sent - " + message)

	// buffer to get data
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	log.Println("Server: message receive - " + string(received))

	return string(received)
}
