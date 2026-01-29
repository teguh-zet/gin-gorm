package loans

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type LoanNatsServer interface {
	Init(svr chan string)
}

type loanNatsServer struct {
	database *gorm.DB
	nc       *nats.Conn
}

func NewLoanNatsServer(db *gorm.DB, nc *nats.Conn, autoMigrate bool) LoanNatsServer {
	if autoMigrate {
		if err := db.AutoMigrate(&LoanLog{}); err != nil {
			log.Printf("Failed to auto migrate LoanLog: %v", err)
		} else {
			log.Println("AutoMigrate LoanLog success")
		}
	}
	return &loanNatsServer{
		database: db,
		nc:       nc,
	}
}

func (s *loanNatsServer) Init(svr chan string) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Panic: %v", r)
		}
		svr <- "loans"
	}()

	loanNatsService := NewLoanNatsService(s.database, s.nc)
	loanNatsControl := NewLoanNatsController(loanNatsService)

	// Topics to listen to
	topics := []string{"book.borrowed", "book_returned"}
	ch := make(chan *nats.Msg, 1024)

	for _, topic := range topics {
		sub, err := s.nc.ChanSubscribe(topic, ch)
		if err != nil {
			log.Printf("can't subscribe to %s: %v", topic, err)
			continue
		}
		fmt.Println("loans -> listen to topic : " + topic)
		defer sub.Unsubscribe()
	}

	// Loop process messages
	for msg := range ch {

		if err := msg.Ack(); err != nil {
			// Just debug log, as it is expected to fail on standard NATS
			log.Printf("Debug: msg.Ack failed (expected for standard NATS): %v", err)
		}

		// Handle Request-Reply pattern (Reply to publisher)
		if msg.Reply != "" {
			s.nc.Publish(msg.Reply, []byte("OK"))
		}

		processMsg(msg, loanNatsControl)
	}

	close(ch)
	svr <- "loans"
}

func processMsg(m *nats.Msg, ctrl LoanNatsController) {
	// pl, ok := parsePayload(m.Data)
	// if !ok {
	// 	log.Printf("Payload tidak valid: %s", string(m.Data))
	// 	return
	// }
var payload PayloadLoan
    if err := json.Unmarshal(m.Data, &payload); err != nil {
        log.Printf("Gagal parsing JSON: %v", err)
        return
    }
	if payload.BookID == 0 || payload.UserID == 0 {
        log.Printf("Data tidak lengkap: %+v", payload)
        return
    }
	switch m.Subject {
	case "book.borrowed":
		ctrl.ProcessBorrow(payload)
	case "book_returned":
		ctrl.ProcessReturn(payload)
	default:
		log.Printf("Unknown subject: %s", m.Subject)
	}
}

