package models

import (
	"time"
)

type Loan struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`                        // ID Peminjam
	User       User      `json:"user" gorm:"foreignKey:UserID"`  // Relasi ke User
	BookID     uint      `json:"book_id"`                        // ID Buku
	Book       Book      `json:"book" gorm:"foreignKey:BookID"`  // Relasi ke Book
	LoanDate   time.Time `json:"loan_date"`                      // Tanggal Pinjam
	ReturnDate time.Time `json:"return_date"`                    // Tanggal Harus Kembali
	Status     string    `json:"status" gorm:"default:borrowed"` // Status: borrowed/returned
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Loan) TableName() string {
	return "loans"
}

type LoanRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}