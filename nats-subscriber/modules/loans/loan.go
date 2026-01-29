package loans

import "time"

type PayloadLoan struct {
	BookID uint `json:"book_id"`
	UserID uint `json:"user_id"`
}

type LoanLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BookID    uint      `json:"book_id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"` // "BORROW" or "RETURN"
	CreatedAt time.Time `json:"created_at"`
}
