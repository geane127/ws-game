package websocket

import "fmt"

type Pool struct {
	Register       chan *Client
	Unregister     chan *Client
	Clients        map[*Client]bool
	Broadcast      chan WinResMessage
	BroadcastStart chan GameStartMessage
	DirectReg      chan RegMsgClient
	DirectGuess    chan GuessMsgClient
	DirectError    chan ErrorMsgClient
}

func NewPool() *Pool {
	return &Pool{
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Clients:        make(map[*Client]bool),
		Broadcast:      make(chan WinResMessage),
		BroadcastStart: make(chan GameStartMessage),
		DirectReg:      make(chan RegMsgClient),
		DirectGuess:    make(chan GuessMsgClient),
		DirectError:    make(chan ErrorMsgClient),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("New User connected")
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// for client, _ := range pool.Clients {
			// 	fmt.Println(client)
			// client.Conn.WriteJSON(Msg{Type: 1, Body: "New User Joined..."})
			// }
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("User disconnected : ", client.PlayerName)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			// for client, _ := range pool.Clients {
			// client.Conn.WriteJSON(Msg{Type: 1, Body: "User Disconnected..."})
			// }
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		case message := <-pool.BroadcastStart:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		case directmessage := <-pool.DirectReg:
			fmt.Println("Sending message to one client : ", directmessage.client.PlayerName)
			if err := directmessage.client.Conn.WriteJSON(directmessage.message); err != nil {
				fmt.Println(err)
				return
			}
		case directmessage := <-pool.DirectGuess:
			fmt.Println("Sending message to one client : ", directmessage.client.PlayerName)
			if err := directmessage.client.Conn.WriteJSON(directmessage.message); err != nil {
				fmt.Println(err)
				return
			}
		case directmessage := <-pool.DirectError:
			fmt.Println("Sending message to one client : ", directmessage.client.PlayerName)
			if err := directmessage.client.Conn.WriteJSON(directmessage.message); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
