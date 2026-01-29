package helper

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn

func ConnectNats(url string) {
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url,
		nats.Timeout(10*time.Second),
		nats.ReconnectWait(5*time.Second),
		nats.MaxReconnects(100),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Disconnected from NATS: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Println("Reconnected to NATS")
		}),
	)

	if err != nil {
		log.Fatal("failed to connect to nats")
		return
	}

	NatsConn = nc
	log.Println("✅ Terhubung ke NATS!")
}
