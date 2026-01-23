package repositories

import (
	"gin-gonic/models"
	"gorm.io/gorm"
)

type BookRepository interface {
	FindAll(offset int, limit int, sortBy string, order string) ([]models.Book, int64, error)
	FindByID(id uint) (models.Book, error)
	Create(book models.Book) (models.Book, error)
	Update(book models.Book) (models.Book, error)
	Delete(book models.Book) error
	UpdateStock(bookID uint, quantity int) error // Helper khusus stok
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (r *bookRepository) FindAll(offset int, limit int, sortBy string, order string) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	// Hitung total dulu
	r.db.Model(&models.Book{}).Count(&total)

	// Ambil data
	err := r.db.Order(sortBy + " " + order).
		Offset(offset).
		Limit(limit).
		Find(&books).Error

	return books, total, err
}

func (r *bookRepository) FindByID(id uint) (models.Book, error) {
	var book models.Book
	err := r.db.First(&book, id).Error
	return book, err
}

func (r *bookRepository) Create(book models.Book) (models.Book, error) {
	err := r.db.Create(&book).Error
	return book, err
}

func (r *bookRepository) Update(book models.Book) (models.Book, error) {
	err := r.db.Save(&book).Error
	return book, err
}

func (r *bookRepository) Delete(book models.Book) error {
	return r.db.Delete(&book).Error
}

func (r *bookRepository) UpdateStock(bookID uint, change int) error {
	// gorm.Expr digunakan agar aman dari race condition (stock = stock + change)
	return r.db.Model(&models.Book{}).Where("id = ?", bookID).
		Update("stock", gorm.Expr("stock + ?", change)).Error
}