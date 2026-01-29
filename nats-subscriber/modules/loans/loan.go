package loans

import "time"

type PayloadLoan struct {
	BookID uint `json:"book_id"`
	UserID uint `json:"user_id"`
	LoanID uint `json:"loan_id"`
}

type LoanLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	LoanID    uint      `json:"loan_id"`
	BookID    uint      `json:"book_id"`
	UserID    uint      `json:"user_id"`
	Action    string    `json:"action"` // "BORROW" or "RETURN"
	CreatedAt time.Time `json:"created_at"`
}
