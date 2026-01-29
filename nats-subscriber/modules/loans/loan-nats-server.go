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
		// ACK is not really needed for Standard NATS Request-Reply,
		// but we include it if user insists on "example code" style.
		// Standard NATS msg.Ack() might return error, so we log it but don't fail.
		if err := msg.Ack(); err != nil {
			// Just debug log, as it is expected to fail on standard NATS
			// log.Printf("Debug: msg.Ack failed (expected for standard NATS): %v", err)
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

// func parsePayload(data []byte) (PayloadLoan, bool) {
// 	pl := PayloadLoan{}
// 	if err := json.Unmarshal(data, &pl); err == nil {
// 		if pl.BookID != 0 || pl.UserID != 0 {
// 			return pl, true
// 		}
// 	}

// 	var env struct {
// 		Data json.RawMessage `json:"data"`
// 	}
// 	if err := json.Unmarshal(data, &env); err == nil && len(env.Data) > 0 {
// 		if err := json.Unmarshal(env.Data, &pl); err == nil {
// 			if pl.BookID != 0 || pl.UserID != 0 {
// 				return pl, true
// 			}
// 		}

// 		var m map[string]interface{}
// 		if err := json.Unmarshal(env.Data, &m); err == nil {
// 			if v, ok := m["book_id"]; ok {
// 				pl.BookID = toUint(v)
// 			}
// 			if v, ok := m["user_id"]; ok {
// 				pl.UserID = toUint(v)
// 			}
// 			if pl.BookID != 0 || pl.UserID != 0 {
// 				return pl, true
// 			}
// 		}
// 	}

// 	var m map[string]interface{}
// 	if err := json.Unmarshal(data, &m); err == nil {
// 		if v, ok := m["book_id"]; ok {
// 			pl.BookID = toUint(v)
// 		}
// 		if v, ok := m["user_id"]; ok {
// 			pl.UserID = toUint(v)
// 		}
// 		if pl.BookID != 0 || pl.UserID != 0 {
// 			return pl, true
// 		}
// 	}

// 	return PayloadLoan{}, false
// }

// func toUint(v interface{}) uint {
// 	switch t := v.(type) {
// 	case float64:
// 		if t < 0 {
// 			return 0
// 		}
// 		return uint(t)
// 	case int:
// 		if t < 0 {
// 			return 0
// 		}
// 		return uint(t)
// 	case int64:
// 		if t < 0 {
// 			return 0
// 		}
// 		return uint(t)
// 	case uint:
// 		return t
// 	case uint64:
// 		return uint(t)
// 	default:
// 		return 0
// 	}
// }
