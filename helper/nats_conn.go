package helper

import (
	"github.com/nats-io/nats.go"
	"log"
)

var NatsConn *nats.Conn
var NatsJS nats.JetStreamContext // Interface untuk JetStream

func ConnectNats(url string) {
	if url == "" {
		url = nats.DefaultURL
	}
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal("failed to connect to nats")
		return
	}

	// Inisiasi JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("GAGAL INIT JETSTREAM", err)
	}

	// Konfigurasi Stream yang diinginkan
	streamName := "LIBRARY_STREAM"
	streamConfig := &nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{"book.*"}, // Wildcard: Menangkap book.borrowed, book.returned, book.stats
		Storage:  nats.FileStorage,
	}

	// Coba buat stream baru
	_, err = js.AddStream(streamConfig)
	if err != nil {
		// Jika error karena stream sudah ada, kita LAKUKAN UPDATE
		// Ini penting agar subjek baru (book.returned) dikenali
		if _, errUpdate := js.UpdateStream(streamConfig); errUpdate != nil {
			log.Printf("⚠️ Gagal Update Stream: %v", errUpdate)
		} else {
			log.Println("✅ Stream berhasil di-update dengan konfigurasi terbaru (book.*)!")
		}
	} else {
		log.Println("✅ Stream baru berhasil dibuat!")
	}

	NatsConn = nc
	NatsJS = js
	log.Println("✅ Terhubung ke NATS JetStream!")
}