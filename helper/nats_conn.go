package helper

import (
	"github.com/nats-io/nats.go"
	"log"
)

var NatsConn *nats.Conn
var NatsJS nats.JetStreamContext //interface untuk jetstream

func ConnectNats(url string) {
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal("failed to connect to nats")
		return
	}
	// inisiasi JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("GAGAL INIT JETSTREAM", err)
	}
	// buat stream
	streamName := "LIBRARY_STREAM"
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"book.*"}, // menangkap semua topik yang berawalan book
		Storage:  nats.FileStorage,
	})
	if err != nil {
		// Jika error karena stream sudah ada, itu wajar/bagus
		log.Printf("ℹ️ Info Stream: %v", err)
	}
	NatsConn = nc
	NatsJS = js
	log.Println("✅ Terhubung ke NATS JetStream!")

}
