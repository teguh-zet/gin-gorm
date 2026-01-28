package loans

import (
	"encoding/json"
	"errors"
	"fmt"

	// "log"
	"time"

	"gin-gonic/helper"
	"gin-gonic/modules/books"

	"gorm.io/gorm"
)

type LoanStats struct {
	TotalTransactions int64 `json:"total_transactions"`
	CurrentlyBorrowed int64 `json:"currently_borrowed"`
	ReturnedBooks     int64 `json:"returned_books"`
}

type LoanService interface {
	GetStats() (*LoanStats, error)
	GetPopularBooks() ([]books.Book, error)
	Borrow(userID uint, input *LoanRequest) (*Loan, error)
	Return(id string) error
	GetMy(userID uint) ([]Loan, error)
	GetAll() ([]Loan, error)
}

type loanService struct {
	db *gorm.DB
}

func NewLoanService(db *gorm.DB) LoanService {
	return &loanService{db: db}
}

//helper function untuk broadcast statistik (agar tidak duplikasi kode karenak dipakai oleh borrow dan return)

func (s *loanService) GetStats() (*LoanStats, error) {
	var totalLoans int64
	var activeLoans int64
	var returnedLoans int64

	if err := s.db.Model(&Loan{}).Count(&totalLoans).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&Loan{}).Where("status = ?", "borrowed").Count(&activeLoans).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&Loan{}).Where("status = ?", "returned").Count(&returnedLoans).Error; err != nil {
		return nil, err
	}

	return &LoanStats{
		TotalTransactions: totalLoans,
		CurrentlyBorrowed: activeLoans,
		ReturnedBooks:     returnedLoans,
	}, nil
}
func (s*loanService) broadcastStats(){
	//1 ambil statistik terru dati db
	stats,err:=s.GetStats()
	if err!=nil{
		fmt.Printf("Gagal mengambil stats untuk broadcast : %v\n",err)
		return
	}
	statsPayload :=map[string]interface{}{
		"type":"STATS_UPDATE",
		"data":stats,
		"time":time.Now(),
	}
	jsonPayload, _:= json.Marshal(statsPayload)

	if helper.NatsJS!=nil{
		_,err:=helper.NatsJS.Publish("book.stats",jsonPayload)
		if err != nil{
			fmt.Printf("⚠️ Gagal publish stats: %v\n", err)
		} else {
			fmt.Println("📊 Stats terbaru terkirim ke JetStream!")
		}
	}

}
func (s *loanService) GetPopularBooks() ([]books.Book, error) {
	var booksData []books.Book
	if err := s.db.Order("borrow_count DESC").Limit(1).Find(&booksData).Error; err != nil {
		return nil, err
	}
	return booksData, nil
}

// func (s *loanService) Borrow(userID uint, input *LoanRequest) (*Loan, error) {
// 	var book books.Book
// 	if err := s.db.First(&book, input.BookID).Error; err != nil {
// 		return nil, errors.New("book not found")
// 	}
// 	if book.Stock <= 0 {
// 		return nil, errors.New("stok habis")
// 	}

// 	tx := s.db.Begin()

// 	if err := tx.Model(&books.Book{}).Where("id = ?", input.BookID).
// 		Update("stock", gorm.Expr("stock - ?", 1)).Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}

// 	loan := Loan{
// 		UserID:     userID,
// 		BookID:     input.BookID,
// 		LoanDate:   time.Now(),
// 		ReturnDate: time.Now().AddDate(0, 0, 7),
// 		Status:     "borrowed",
// 	}

// 	if err := tx.Create(&loan).Error; err != nil {
// 		tx.Rollback()
// 		return nil, err
// 	}
// 	//mengganti nats yang lebih ringan
// 	// if err := tx.Model(&books.Book{}).Where("id = ?", input.BookID).
// 	// 	Update("borrow_count", gorm.Expr("borrow_count + ?", 1)).Error; err != nil {
// 	// 	tx.Rollback()
// 	// 	return nil, err
// 	// }

// 	if err := tx.Commit().Error; err != nil {
// 		return nil, err
// 	}
// 	// ======================================================
// 	// DEBUGGING AREA - CEK DISINI
// 	// ======================================================
// 	fmt.Println("👉 DEBUG: Transaksi DB Selesai. Mencoba Publish NATS...")

// 	eventData := map[string]interface{}{
// 		"book_id": input.BookID,
// 		"user_id": userID,
// 	}
// 	payload, _ := json.Marshal(eventData)

// 	if helper.NatsConn == nil {
// 		fmt.Println("❌ ERROR: helper.NatsConn is NIL (Kosong) di Loan Service!")
// 	} else {
// 		err := helper.NatsConn.Publish("book.borrowed", payload)
// 		if err != nil {
// 			fmt.Printf("❌ ERROR: Gagal Publish NATS: %v\n", err)
// 		} else {
// 			fmt.Println("✅ SUKSES: Pesan NATS Terkirim ke topik 'book.borrowed'!")
// 		}
// 	}
// 	var fullLoan Loan
// 	if err := s.db.Preload("User").Preload("Book").First(&fullLoan, loan.ID).Error; err != nil {
// 		return nil, err
// 	}

// 	return &fullLoan, nil
// }
func (s *loanService) Borrow(userID uint, input *LoanRequest) (*Loan, error) {
    var book books.Book
    if err := s.db.First(&book, input.BookID).Error; err != nil {
        return nil, errors.New("book not found")
    }
    if book.Stock <= 0 {
        return nil, errors.New("stok habis")
    }

    tx := s.db.Begin()

    if err := tx.Model(&books.Book{}).Where("id = ?", input.BookID).
        Update("stock", gorm.Expr("stock - ?", 1)).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    loan := Loan{
        UserID:     userID,
        BookID:     input.BookID,
        LoanDate:   time.Now(),
        ReturnDate: time.Now().AddDate(0, 0, 7),
        Status:     "borrowed",
    }

    if err := tx.Create(&loan).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

 
	// eventData := map[string]interface{}{
    //     "book_id": input.BookID,
    //     "user_id": userID,
    // }
    // payload, _ := json.Marshal(eventData)


    // if helper.NatsJS != nil {
    //     // Ganti helper.NatsConn menjadi helper.NatsJS
    //     // Tidak perlu Flush() manual, JetStream menghandle ack-nya
    //     _, err := helper.NatsJS.Publish("book.borrowed", payload)
        
    //     if err != nil {
    //         fmt.Printf("⚠️ Gagal publish JetStream: %v\n", err) 
    //     } else {
    //          fmt.Println("✅ SUKSES: Pesan tersimpan di JetStream!")
    //     }
    // }
    // ------------------------------------
	// 1. Buat Payload Data untuk Frontend
    // Structure ini nanti akan diterima oleh WebSocket di browser
    notificationData := map[string]interface{}{
        "type":    "NEW_LOAN", // Penanda tipe event
        "message": fmt.Sprintf("User ID %d baru saja meminjam buku ID %d", userID, input.BookID),
        "data": map[string]interface{}{
            "book_id": input.BookID,
            "user_id": userID,
			"action" : "borrow",
            "time":    time.Now(),
        },
    }

    // 2. Ubah ke JSON
    payload, _ := json.Marshal(notificationData)

    // 3. Publish ke NATS
    // PENTING: Pastikan string "notifications" SAMA PERSIS dengan 
    // yang ada di manager.go (helper.NatsConn.Subscribe("notifications", ...))
    if helper.NatsJS != nil {
        _, err := helper.NatsJS.Publish("book.borrowed", payload)
        
        if err != nil {
            // Log error tapi jangan gagalkan request user karena DB sudah sukses
            fmt.Printf("⚠️ Gagal publish ke JetStream: %v\n", err)
        } else {
            fmt.Println("✅ SUKSES: Notifikasi tersimpan di JetStream & siap dikirim!")
        }
    } else {
        fmt.Println("❌ Error: Koneksi JetStream belum siap (nil)")
    }
	// goroutine agar tidak memblokir response API
	go s.broadcastStats()

    var fullLoan Loan
    if err := s.db.Preload("User").Preload("Book").First(&fullLoan, loan.ID).Error; err != nil {
        return nil, err
    }

    return &fullLoan, nil
}
func (s *loanService) Return(id string) error {
	var loan Loan
	if err := s.db.First(&loan, id).Error; err != nil {
		return errors.New("loan not found")
	}
	if loan.Status == "returned" {
		return errors.New("book already returned")
	}

	tx := s.db.Begin()

	if err := tx.Model(&Loan{}).Where("id = ?", loan.ID).
		Updates(map[string]interface{}{"status": "returned", "return_date": time.Now()}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&books.Book{}).Where("id = ?", loan.BookID).
		Update("stock", gorm.Expr("stock + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error;err!=nil{
		return err
	}
	// nats implementation for returning notification
	eventData:= map[string]interface{}{
		"type" : "BOOK_RETURNED",
		"message":fmt.Sprintf("Buku (loan ID : %s)telah dikembalikan", id),
		"data":map[string]interface{}{
			"loan_id":id,
			"book_id": loan.BookID,
			"action":"return",
			"time":time.Now(),

		},
	}
	payload,_:= json.Marshal(eventData)

	if helper.NatsJS != nil{
		_,err := helper.NatsJS.Publish("book_returned", payload)
		if err!=nil{
			fmt.Printf("GAGAL PUBLISH RETURN EVENT: %v\n",err)
		}else{
			fmt.Println("NOTIFIKASI RETURN TERSIMPAN DI JETSTREAM")
		}
	}
	go s.broadcastStats()

	return nil
}

func (s *loanService) GetMy(userID uint) ([]Loan, error) {
	var loansData []Loan
	if err := s.db.Preload("Book").Where("user_id = ?", userID).Find(&loansData).Error; err != nil {
		return nil, err
	}
	return loansData, nil
}

func (s *loanService) GetAll() ([]Loan, error) {
	var loansData []Loan
	if err := s.db.Preload("User").Preload("Book").Find(&loansData).Error; err != nil {
		return nil, err
	}
	return loansData, nil
}
