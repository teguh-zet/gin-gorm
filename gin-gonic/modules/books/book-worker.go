package books

import (
	"encoding/json"
	"gin-gonic/helper"
	"log"

	"github.com/nats-io/nats.go"
)

type BookBorrowedEvent struct {
	BookID uint `json:"book_id"`
	UserID uint `json:"user_id"`
}

func StartWorker(service BookService) {
	if helper.NatsConn == nil {
		log.Println("⚠️ Worker Batal: NATS Conn is NIL")
		return
	}

	// Gunakan Subscribe biasa
	_, err := helper.NatsConn.Subscribe("book.borrowed", func(msg *nats.Msg) {

		// 1. Parsing Data
		var event BookBorrowedEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("❌ Gagal parsing: %v", err)
			return
		}

		log.Printf("📩 [NATS] Update popularitas Buku ID: %d", event.BookID)

		// 2. Eksekusi Logic DB
		err = service.IncrementPopularity(event.BookID)
		if err != nil {
			log.Printf("❌ DB Error: %v", err)
			return
		}

		// Reply to request if needed
		if msg.Reply != "" {
			helper.NatsConn.Publish(msg.Reply, []byte("OK"))
		}

		log.Println("✅ Proses Selesai")

	})

	if err != nil {
		log.Fatal("❌ Gagal Subscribe:", err)
	}

	log.Println("🎧 NATS Worker siap mendengarkan...")
}
