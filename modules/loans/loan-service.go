package loans

import (
	"net/http"
	"time"

	"gin-gonic/helper"
	"gin-gonic/modules/books"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoanService interface {
	GetStats(ctx *gin.Context)
	GetPopularBooks(ctx *gin.Context)
	Borrow(ctx *gin.Context)
	Return(ctx *gin.Context)
	GetMy(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}

type loanService struct {
	db *gorm.DB
}

func NewLoanService(db *gorm.DB) LoanService {
	return &loanService{db: db}
}

func (s *loanService) GetStats(c *gin.Context) {
	var totalLoans int64
	var activeLoans int64
	var returnedLoans int64

	s.db.Model(&Loan{}).Count(&totalLoans)
	s.db.Model(&Loan{}).Where("status = ?", "borrowed").Count(&activeLoans)
	s.db.Model(&Loan{}).Where("status = ?", "returned").Count(&returnedLoans)

	stats := gin.H{
		"total_transactions": totalLoans,
		"currently_borrowed": activeLoans,
		"returned_books":     returnedLoans,
	}
	helper.SuccessResponse(c, "Loan statistics", stats)
}

func (s *loanService) GetPopularBooks(c *gin.Context) {
	var books []books.Book
	s.db.Order("borrow_count DESC").Limit(1).Find(&books)
	helper.SuccessResponse(c, "Popular books", books)
}

func (s *loanService) Borrow(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userIDFloat, ok := userIDVal.(float64)
	if !ok {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}
	userID := uint(userIDFloat)

	var req LoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationError(c, err.Error())
		return
	}

	var book books.Book
	if err := s.db.First(&book, req.BookID).Error; err != nil {
		helper.NotFoundError(c, "Book not found")
		return
	}

	if book.Stock <= 0 {
		helper.ErrorResponse(c, http.StatusBadRequest, "Stok habis", "Buku tidak tersedia")
		return
	}

	tx := s.db.Begin()

	if err := tx.Model(&books.Book{}).Where("id = ?", req.BookID).
		Update("stock", gorm.Expr("stock - ?", 1)).Error; err != nil {
		tx.Rollback()
		helper.InternalServerError(c, "Failed to update stock", err.Error())
		return
	}

	loan := Loan{
		UserID:     userID,
		BookID:     req.BookID,
		LoanDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 7),
		Status:     "borrowed",
	}

	if err := tx.Create(&loan).Error; err != nil {
		tx.Rollback()
		helper.InternalServerError(c, "Failed to create loan", err.Error())
		return
	}

	if err := tx.Model(&books.Book{}).Where("id = ?", req.BookID).
		Update("borrow_count", gorm.Expr("borrow_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		helper.InternalServerError(c, "Failed to update borrow count", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		helper.InternalServerError(c, "Failed to commit transaction", err.Error())
		return
	}

	var fullLoan Loan
	if err := s.db.Preload("User").Preload("Book").First(&fullLoan, loan.ID).Error; err != nil {
		helper.InternalServerError(c, "Failed to fetch loan", err.Error())
		return
	}

	helper.SuccessResponse(c, "Book borrowed successfully", fullLoan)
}

func (s *loanService) Return(c *gin.Context) {
	id := c.Param("id")

	var loan Loan
	if err := s.db.First(&loan, id).Error; err != nil {
		helper.NotFoundError(c, "Loan not found")
		return
	}

	if loan.Status == "returned" {
		helper.ErrorResponse(c, http.StatusBadRequest, "Buku sudah dikembalikan", nil)
		return
	}

	tx := s.db.Begin()

	if err := tx.Model(&Loan{}).Where("id = ?", loan.ID).
		Updates(map[string]interface{}{"status": "returned", "return_date": time.Now()}).Error; err != nil {
		tx.Rollback()
		helper.InternalServerError(c, "Failed to update loan", err.Error())
		return
	}

	if err := tx.Model(&books.Book{}).Where("id = ?", loan.BookID).
		Update("stock", gorm.Expr("stock + ?", 1)).Error; err != nil {
		tx.Rollback()
		helper.InternalServerError(c, "Failed to update stock", err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		helper.InternalServerError(c, "Failed to commit transaction", err.Error())
		return
	}

	helper.SuccessResponse(c, "Book returned successfully", nil)
}

func (s *loanService) GetMy(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userIDFloat, ok := userIDVal.(float64)
	if !ok {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}
	userID := uint(userIDFloat)

	var loans []Loan
	if err := s.db.Preload("Book").Where("user_id = ?", userID).Find(&loans).Error; err != nil {
		helper.InternalServerError(c, "Failed to fetch loans", err.Error())
		return
	}

	helper.SuccessResponse(c, "My loans", loans)
}

func (s *loanService) GetAll(c *gin.Context) {
	var loans []Loan
	if err := s.db.Preload("User").Preload("Book").Find(&loans).Error; err != nil {
		helper.InternalServerError(c, "Failed to fetch loans", err.Error())
		return
	}

	helper.SuccessResponse(c, "All loans", loans)
}
