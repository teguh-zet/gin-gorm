package websocket

import (
	"log"
	"sync"
)

type Manager struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Mutex      sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (m *Manager) Run() {
	// =================================================================
	// SUBSCRIBER REMOVED - MOVED TO EXTERNAL SERVICE
	// =================================================================

	for {
		select {
		case client := <-m.Register:
			m.Mutex.Lock()
			m.Clients[client] = true
			m.Mutex.Unlock()
			log.Printf("Client Connected: %s", client.UserID)

		case client := <-m.Unregister:
			m.Mutex.Lock()
			if _, ok := m.Clients[client]; ok {
				delete(m.Clients, client)
				close(client.Send)
			}
			m.Mutex.Unlock()
			log.Printf("Client Disconnected: %s", client.UserID)

		case message := <-m.Broadcast:
			m.Mutex.Lock()
			for client := range m.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.Clients, client)
				}
			}
			m.Mutex.Unlock()
		}
	}
}
