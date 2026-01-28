package websocket

import (
	
	"gin-gonic/helper"
	"log"
	"sync"

	// Sesuaikan import path ini dengan module name di go.mod Anda

	"github.com/nats-io/nats.go"
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
    // 1. Subscribe ke NATS (Topic: "notifications")
    // Ini menghubungkan NATS dengan WebSocket
	_, err := helper.NatsConn.Subscribe("notifications", func(msg *nats.Msg) {
		m.Broadcast <- msg.Data
	})
    if err != nil {
        log.Printf("Error subscribing to NATS: %v", err)
    }

	for {
		select {
		case client := <-m.Register:
			m.Mutex.Lock()
			m.Clients[client] = true
			m.Mutex.Unlock()
			log.Printf("New client connected. UserID: %s", client.UserID)

		case client := <-m.Unregister:
			m.Mutex.Lock()
			if _, ok := m.Clients[client]; ok {
				delete(m.Clients, client)
				close(client.Send)
			}
			m.Mutex.Unlock()
			log.Printf("Client disconnected. UserID: %s", client.UserID)

		case message := <-m.Broadcast:
            // Mengirim pesan ke SEMUA client yang terhubung
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