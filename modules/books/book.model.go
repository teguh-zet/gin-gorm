package books

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string `json:"title" gorm:"not null"`
	Author      string `json:"author"`
	Stock       int    `json:"stock" gorm:"default:0"` // [NEW] Menyimpan jumlah stok
	BorrowCount int    `json:"borrow_count" gorm:"default:0"`
	ImageURL    string `json:"image_url"`
	// fine        int64          `json:"fine" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`     // [NEW] Waktu dibuat
	UpdatedAt time.Time      `json:"updated_at"`     // [NEW] Waktu terakhir diedit
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // [NEW] Soft delete
}

func (Book) TableName() string {
	return "books"
}

// Struct untuk validasi saat membuat buku baru
type CreateBookRequest struct {
	Title       string `json:"title" binding:"required,min=2,max=100"`
	Author      string `json:"author" binding:"required,min=2,max=100"`
	Stock       int    `json:"stock" binding:"required,min=0"` // [NEW] Wajib isi stok minimal 0
	BorrowCount int    `json:"borrow_count" gorm:"default:0"`
}

// Struct untuk validasi saat update buku (opsional fieldnya)
type UpdateBookRequest struct {
	Title  string `json:"title" binding:"omitempty,min=2,max=100"`
	Author string `json:"author" binding:"omitempty,min=2,max=100"`
	Stock  int    `json:"stock" binding:"omitempty,min=0"` // [NEW] Update stok
}
