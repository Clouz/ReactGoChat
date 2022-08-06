package websocket

import "fmt"

type Pool struct {
	Register     chan *Client
	Unregister   chan *Client
	Clients      map[*Client]bool
	lastClientID int
	Broadcast    chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Clients:      make(map[*Client]bool),
		Broadcast:    make(chan Message),
		lastClientID: 0,
	}
}

func (pool *Pool) GetClientID() string {
	pool.lastClientID++
	return fmt.Sprintf("User%d", pool.lastClientID)
}

func (pool *Pool) Start() {
	for {
		select {
		case registerClient := <-pool.Register:
			pool.Clients[registerClient] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				fmt.Println(client)
				if registerClient == client {
					client.Conn.WriteJSON(Message{Type: 1, ID: registerClient.ID, Body: "Welcome"})
				} else {
					client.Conn.WriteJSON(Message{Type: 1, ID: registerClient.ID, Body: "New User Joined..."})
				}
			}
			break
		case unregisterClient := <-pool.Unregister:
			delete(pool.Clients, unregisterClient)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, ID: unregisterClient.ID, Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
