package books

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"gin-gonic/helper"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"gorm.io/gorm"
)

type BookService interface {
	GetList() ([]Book, error)
	GetList2(page, limit int, sortBy, order string, available bool) ([]Book, int64, error)
	GetByID(id string) (*Book, error)
	Create(input *CreateBookRequest) (*Book, error)
	Update(id string, input *UpdateBookRequest) (*Book, error)
	Delete(id string) error
	BulkDelete(ids []int) error
	Search(query string) ([]Book, error)
	UploadImage(id string, file *multipart.FileHeader) (*Book, error)
	IncrementPopularity(bookID uint)error
}

type bookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) BookService {
	return &bookService{db: db}
}
func (s *bookService) IncrementPopularity(bookID uint) error {
	// Debug log
	fmt.Printf("⚙️ [SERVICE] Menjalankan Query Update untuk Book ID: %d\n", bookID)

	// Update borrow_count + 1
	err := s.db.Model(&Book{}).Where("id = ?", bookID).
		Update("borrow_count", gorm.Expr("borrow_count + 1")).Error
	
	if err != nil {
		fmt.Printf("❌ [SERVICE ERROR] Gorm Error: %v\n", err)
		return err
	}
	
	return nil
}

func (s *bookService) GetList() ([]Book, error) {
	var books []Book
	if err := s.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (s *bookService) GetList2(page, limit int, sortBy, order string, available bool) ([]Book, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	allowedSort := map[string]bool{
		"id":     true,
		"title":  true,
		"author": true,
		"stock":  true,
	}
	if !allowedSort[sortBy] {
		sortBy = "id"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	offset := (page - 1) * limit

	query := s.db.Model(&Book{})
	if available {
		query = query.Where("stock > ?", 0)
	}

	var books []Book
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order(sortBy + " " + order).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (s *bookService) GetByID(id string) (*Book, error) {
	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, errors.New("book not found")
	}
	return &book, nil
}

func (s *bookService) Create(input *CreateBookRequest) (*Book, error) {
	book := &Book{
		Title:       input.Title,
		Author:      input.Author,
		Stock:       input.Stock,
		BorrowCount: 0,
	}

	if err := s.db.Create(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) Update(id string, input *UpdateBookRequest) (*Book, error) {
	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, errors.New("book not found")
	}

	if input.Title != "" {
		book.Title = input.Title
	}
	if input.Author != "" {
		book.Author = input.Author
	}
	if input.Stock != 0 {
		book.Stock = input.Stock
	}

	if err := s.db.Save(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *bookService) Delete(id string) error {
	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		return errors.New("book not found")
	}

	return s.db.Delete(&book).Error
}

func (s *bookService) BulkDelete(ids []int) error {
	return s.db.Delete(&[]Book{}, ids).Error
}

func (s *bookService) Search(query string) ([]Book, error) {
	if query == "" {
		return nil, errors.New("query parameter q is required")
	}

	var books []Book
	if err := s.db.Where("title ILIKE ? OR author ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (s *bookService) UploadImage(id string, file *multipart.FileHeader) (*Book, error) {
	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, errors.New("book not found")
	}

	config, err := helper.LoadConfig(".")
	if err != nil {
		return nil, errors.New("failed to load config")
	}

	if config.CloudinaryCloudName == "" || config.CloudinaryAPIKey == "" || config.CloudinaryAPISecret == "" {
		return nil, errors.New("cloudinary config missing")
	}

	cld, err := cloudinary.NewFromParams(
		config.CloudinaryCloudName,
		config.CloudinaryAPIKey,
		config.CloudinaryAPISecret,
	)
	if err != nil {
		return nil, err
	}

	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	uploadResult, err := cld.Upload.Upload(context.Background(), fileReader, uploader.UploadParams{})
	if err != nil {
		return nil, err
	}

	book.ImageURL = uploadResult.SecureURL
	if err := s.db.Save(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}
