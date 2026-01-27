package books

import (
	"encoding/json"
	"gin-gonic/helper" // Sesuaikan dengan nama modul Anda (gin-gonic atau gin-gorm)
	"log"

	"github.com/nats-io/nats.go"
)

type BookBorrowedEvent struct {
	BookID uint `json:"book_id"`
	UserID uint `json:"user_id"`
}

func StartWorker(service BookService) {
	if helper.NatsJS == nil {
		log.Println("⚠️ Worker Batal: NATS JetStream is NIL")
		return
	}

	// Gunakan Subscribe biasa (Push Consumer) via JetStream
	_, err := helper.NatsJS.Subscribe("book.borrowed", func(msg *nats.Msg) {
		
        // 1. Parsing Data
		var event BookBorrowedEvent
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("❌ Gagal parsing: %v", err)
            msg.Term() // Terminate: Bilang ke server "Pesan ini rusak, jangan kirim lagi"
			return
		}

		log.Printf("📩 [JETSTREAM] Update popularitas Buku ID: %d", event.BookID)
		
		// 2. Eksekusi Logic DB
		err = service.IncrementPopularity(event.BookID)
		if err != nil {
            // Jika DB error, kita JANGAN Ack. 
            // Nanti NATS akan mengirim ulang pesan ini (Retry) otomatis.
			log.Printf("❌ DB Error: %v", err)
            msg.Nak() // Negative Ack: "Coba kirim lagi nanti"
            return
		}

        // 3. [PENTING] Acknowledge (Tanda Sukses)
        // Kalau ini lupa, pesan akan dikirim ulang terus menerus!
        msg.Ack()
        log.Println("✅ Proses Selesai & Terkonfirmasi (Ack)")

	}, nats.Durable("book_worker_durable"), nats.ManualAck()) // Config agar worker mengingat posisi terakhir

	if err != nil {
		log.Fatal("❌ Gagal Subscribe JS:", err)
	}

	log.Println("🎧 JetStream Worker siap mendengarkan...")
}