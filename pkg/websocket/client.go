package websocket

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

var ran int
var gameid int64 = 1

//Client : Object of the connected client
type Client struct {
	ID         string
	Conn       *websocket.Conn
	Pool       *Pool
	PlayerName string
}

//Msg : General Mesage type in json format
type Msg struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

//RegMsgClient : Registration Message response to specific client
type RegMsgClient struct {
	message RegResMessage
	client  *Client
}

//GuessMsgClient : Guess Message response to specific client
type GuessMsgClient struct {
	message GuessResMessage
	client  *Client
}

//ErrorMsgClient : Error Message response to specific client
type ErrorMsgClient struct {
	message ErrorMessage
	client  *Client
}

//ReqMessage : All type of Request Message structure in json format
type ReqMessage struct {
	Message    string `json:"message"`
	PlayerName string `json:"playerName"`
	Guess      int    `json:"guess"`
	Timestamp  int64  `json:"timestamp"`
	GameID     int64  `json:"gameId"`
}

//RegResMessage : Registration Response Message structure in json format
type RegResMessage struct {
	Message    string `json:"message"`
	PlayerName string `json:"playerName"`
	Timestamp  int64  `json:"timestamp"`
	GameID     int64  `json:"gameId"`
}

//GuessResMessage : Guess Response Message structure in json format
type GuessResMessage struct {
	Message     string `json:"message"`
	GuessResult int    `json:"guessResult"`
	Timestamp   int64  `json:"timestamp"`
	GameID      int64  `json:"gameId"`
}

//WinResMessage : Game Win Message structure in json format
type WinResMessage struct {
	Message string `json:"message"`
	Answer  int    `json:"answer"`
	Winner  string `json:"winner"`
	GameID  int64  `json:"gameId"`
}

//GameStartMessage : Game Start Message structure in json format
type GameStartMessage struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	GameID    int64  `json:"gameId"`
}

//ErrorMessage : Error Message structure in json format
type ErrorMessage struct {
	Message   string `json:"message"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

func (c *Client) Read() {

	var x int

	// Update the random seed
	rand.Seed(time.Now().UTC().UnixNano())
	ran = rand.Intn(499) + 1

	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var reqMessage ReqMessage
		json.Unmarshal([]byte(string(p)), &reqMessage)
		fmt.Println()
		fmt.Println("Message: ", reqMessage.Message, ", PlayerName: ", reqMessage.PlayerName, ", Guess: ", reqMessage.Guess, ", Timestamp: ", reqMessage.Timestamp, ", GameID: ", reqMessage.GameID)

		// Hash value for the word "registration"
		bv := []byte("registration")
		hasher := sha256.New()
		hasher.Write(bv)
		regsha := hex.EncodeToString(hasher.Sum(nil))

		// Hash value for the word "guess"
		bv = []byte("guess")
		hasher = sha256.New()
		hasher.Write(bv)
		guesssha := hex.EncodeToString(hasher.Sum(nil))

		// Display the timestamp from request message
		tm := time.Unix(reqMessage.Timestamp, 0)
		fmt.Println(tm)

		if regsha == reqMessage.Message {
			fmt.Println("Registration: ", regsha)
			c.PlayerName = reqMessage.PlayerName
			reqMessage := RegResMessage{Message: reqMessage.Message, PlayerName: reqMessage.PlayerName, Timestamp: reqMessage.Timestamp, GameID: gameid}
			guess := RegMsgClient{reqMessage, c}
			c.Pool.DirectReg <- guess
		} else if guesssha == reqMessage.Message {
			fmt.Println("Guess: ", guesssha)

			if reqMessage.GameID != gameid {
				errorMessage := ErrorMessage{Message: "ca00fccfb408989eddc401062c4d1219a6aceb6b9b55412357f1790862e8f178", Reason: "696e636f72726563742067616d65204944", Timestamp: int64(time.Now().Unix())}
				error := ErrorMsgClient{errorMessage, c}
				c.Pool.DirectError <- error
			} else {
				x, _ = checkRan(reqMessage.Guess)

				if x == 0 {
					winMessage := WinResMessage{Message: reqMessage.Message, Answer: reqMessage.Guess, Winner: c.PlayerName, GameID: gameid}
					c.Pool.Broadcast <- winMessage
					gameid++
					gameMessage := GameStartMessage{Message: "9b7db65d5e1f739da360fb8c65879114d63d44297fbd274a9e32c05887b5ba99", Timestamp: int64(time.Now().Unix()), GameID: gameid}
					c.Pool.BroadcastStart <- gameMessage
				} else {
					guessMessage := GuessResMessage{Message: reqMessage.Message, GuessResult: x, Timestamp: reqMessage.Timestamp, GameID: gameid}
					guess := GuessMsgClient{guessMessage, c}
					c.Pool.DirectGuess <- guess
				}
			}
		} else {
			fmt.Println("No match")
		}

	}
}

func checkRan(p int) (x int, result string) {

	fmt.Println("Random Number = ", ran)
	fmt.Println("Guess = ", p)

	if p == ran {
		fmt.Println("Bingo")
		result = "Bingo"
		x = 0
		// Update random seed for generating the next random number
		rand.Seed(time.Now().UTC().UnixNano())
		ran = rand.Intn(499) + 1

	} else {
		if p < ran {
			fmt.Println("too small")
			x = 2
			result = "too small"
		} else {
			fmt.Println("too large")
			x = 1
			result = "too large"
		}
	}
	return
}
