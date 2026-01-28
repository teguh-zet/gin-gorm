package websocket

import (
	
	"gin-gonic/helper"
	"log"
	"sync"
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
	// Pastikan koneksi JetStream tersedia
	if helper.NatsJS == nil {
		log.Fatal("❌ ERROR: Koneksi NATS JetStream (NatsJS) belum siap/nil!")
	}

	// =================================================================
	// UPDATE SUBSCRIBER KE JETSTREAM (Durable Consumer)
	// =================================================================
	// 1. Subscribe ke Topic "book.borrowed" menggunakan NatsJS (bukan NatsConn)
	sub, err := helper.NatsJS.Subscribe("book.*", func(msg *nats.Msg) {
		
		// 2. ACKNOWLEDGEMENT (Konfirmasi Terima)
		// Kita wajib lapor ke NATS bahwa pesan sudah diterima.
		// Jika tidak di-Ack, NATS akan menganggap pesan gagal dan mengirim ulang terus.
		if err := msg.Ack(); err != nil {
			log.Printf("⚠️ Gagal Ack pesan: %v", err)
			return // Jangan diproses kalau gagal Ack (opsional, tergantung strategi)
		}

		// 3. Log Debug (Agar Anda tahu pesan masuk)
		log.Printf("📨 [WS Broadcast] Topic: %s | Data: %s", msg.Subject, string(msg.Data))

		// 4. Broadcast ke Frontend
		m.Broadcast <- msg.Data

	}, 
	// OPSI PENTING:
	// nats.Durable: Membuat NATS "mengingat" posisi terakhir kita.
	// Jika server mati lalu nyala, pesan yang terlewat akan dikirim ulang.
	nats.Durable("WS_MANAGER_CONSUMER"), 
	// nats.ManualAck: Kita janji akan melakukan msg.Ack() sendiri secara manual.
	nats.ManualAck(),
	)

	if err != nil {
		log.Printf("❌ Gagal Subscribe JetStream: %v", err)
	} else {
		log.Println("🚀 WebSocket Manager mendengarkan JetStream (Topic: book.borrowed)")
		defer sub.Unsubscribe()
	}
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