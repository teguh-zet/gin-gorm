package loans

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type LoanNatsService interface {
	ProcessBorrow(payload PayloadLoan)
	ProcessReturn(payload PayloadLoan)
}

type loanNatsService struct {
	conn *gorm.DB
	nc   *nats.Conn
}

func NewLoanNatsService(db *gorm.DB, nc *nats.Conn) LoanNatsService {
	return &loanNatsService{
		conn: db,
		nc:   nc,
	}
}

func (service *loanNatsService) ProcessBorrow(payload PayloadLoan) {
	log.Printf("Processing Borrow Event: %+v", payload)

	// Simpan Log ke DB (sesuai contoh rawat_inap yang save ke DB)
	logEntry := LoanLog{
		BookID:    payload.BookID,
		UserID:    payload.UserID,
		Action:    "BORROW",
		CreatedAt: time.Now(),
	}

	// Pastikan table ada (AutoMigrate biasanya di main, tapi kita handle aman di sini atau assume ada)
	// service.conn.AutoMigrate(&LoanLog{}) 

	if err := service.conn.Create(&logEntry).Error; err != nil {
		log.Printf("Gagal menyimpan log borrow ke database: %v", err)
	} else {
		log.Printf("Log borrow berhasil disimpan: %+v", logEntry)
	}
}

func (service *loanNatsService) ProcessReturn(payload PayloadLoan) {
	log.Printf("Processing Return Event: %+v", payload)

	logEntry := LoanLog{
		
		BookID:    payload.BookID,
		UserID:    payload.UserID,
		Action:    "RETURN",
		CreatedAt: time.Now(),
	}

	if err := service.conn.Create(&logEntry).Error; err != nil {
		log.Printf("Gagal menyimpan log return ke database: %v", err)
	} else {
		log.Printf("Log return berhasil disimpan: %+v", logEntry)
	}
}
