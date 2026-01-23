// 

package loan_controller

import (
	"net/http"
	"time"

	"gin-gonic/database"
	"gin-gonic/helpers"
	"gin-gonic/models"

	"github.com/gin-gonic/gin"
)

// BorrowBook godoc
// @Summary      Pinjam Buku
// @Description  Membuat data peminjaman baru dan mengurangi stok buku.
// @Tags         loans
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        request body models.LoanRequest true "Data Peminjaman (Book ID)"
// @Success      201  {object} models.Loan
// @Failure      400  {object} map[string]interface{} "Stok habis atau sedang meminjam"
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Failure      404  {object} map[string]interface{} "Buku tidak ditemukan"
// @Failure      500  {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /loans [post]
func BorrowBook(c *gin.Context) {
	// Ambil User ID dari JWT
	userIDVal, exists := c.Get("user_id")
	if !exists {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	userID := uint(userIDVal.(float64))

	// Pastikan struct LoanRequest ada di models/loan.model.go
	// type LoanRequest struct { BookID uint `json:"book_id"` }
	var req models.LoanRequest 
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err.Error())
		return
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var book models.Book
	if err := tx.First(&book, req.BookID).Error; err != nil {
		tx.Rollback()
		helpers.NotFoundError(c, "Book not found")
		return
	}

	if book.Stock < 1 {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusBadRequest, "Out of stock", "Stok buku habis")
		return
	}

	var activeLoan int64
	tx.Model(&models.Loan{}).
		Where("user_id = ? AND book_id = ? AND status = ?", userID, req.BookID, "borrowed").
		Count(&activeLoan)

	if activeLoan > 0 {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusBadRequest, "Duplicate loan", "Anda sedang meminjam buku ini")
		return
	}

	book.Stock = book.Stock - 1
	if err := tx.Save(&book).Error; err != nil {
		tx.Rollback()
		helpers.InternalServerError(c, "Failed to update stock", err.Error())
		return
	}

	loan := models.Loan{
		UserID:     userID,
		BookID:     req.BookID,
		LoanDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 7),
		Status:     "borrowed",
	}

	if err := tx.Create(&loan).Error; err != nil {
		tx.Rollback()
		helpers.InternalServerError(c, "Failed to create loan", err.Error())
		return
	}

	tx.Commit()

	var fullLoan models.Loan
	if err := database.DB.Preload("User").Preload("Book").First(&fullLoan, loan.ID).Error; err != nil {
		helpers.InternalServerError(c, "Failed to load created loan", err.Error())
		return
	}
	helpers.SuccessResponse(c, "Book borrowed successfully", fullLoan)
}

// ReturnBook godoc
// @Summary      Kembalikan Buku
// @Description  Mengubah status peminjaman menjadi 'returned' dan menambah stok buku. Hanya pemilik peminjaman yang boleh akses.
// @Tags         loans
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        id path int true "Loan ID"
// @Success      200  {object} map[string]interface{} "Pesan Sukses"
// @Failure      400  {object} map[string]interface{} "Sudah dikembalikan"
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Failure      403  {object} map[string]interface{} "Bukan milik user ini"
// @Failure      404  {object} map[string]interface{} "Data Loan tidak ditemukan"
// @Security     BearerAuth
// @Router       /loans/return/{id} [post]
func ReturnBook(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	loanID := c.Param("id")

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var loan models.Loan
	if err := tx.First(&loan, loanID).Error; err != nil {
		tx.Rollback()
		helpers.NotFoundError(c, "Loan record not found")
		return
	}

	if loan.UserID != userID {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusForbidden, "Forbidden", "Ini bukan peminjaman Anda")
		return
	}

	if loan.Status == "returned" {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusBadRequest, "Already returned", "Buku sudah dikembalikan sebelumnya")
		return
	}

	loan.Status = "returned"
	loan.ReturnDate = time.Now()
	if err := tx.Save(&loan).Error; err != nil {
		tx.Rollback()
		helpers.InternalServerError(c, "Failed to update loan", err.Error())
		return
	}

	if err := tx.Model(&models.Book{}).Where("id = ?", loan.BookID).
		Update("stock", database.DB.Raw("stock + 1")).Error; err != nil {
		tx.Rollback()
		helpers.InternalServerError(c, "Failed to update book stock", err.Error())
		return
	}

	tx.Commit()

	helpers.SuccessResponse(c, "Book returned successfully", nil)
}

// GetMyLoans godoc
// @Summary      Lihat Riwayat Peminjaman Saya
// @Description  Menampilkan daftar semua buku yang pernah dipinjam oleh user yang sedang login.
// @Tags         loans
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200  {array} models.Loan
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Security     BearerAuth
// @Router       /loans/my [get]
func GetMyLoans(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	var loans []models.Loan

	if err := database.DB.Preload("Book").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&loans).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch loans", err.Error())
		return
	}

	helpers.SuccessResponse(c, "Your loan history", loans)
}

// GetAllLoans godoc
// @Summary      Lihat Semua Peminjaman (Admin)
// @Description  Menampilkan seluruh data peminjaman dari semua user. Khusus Admin.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200  {array} models.Loan
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Failure      403  {object} map[string]interface{} "Forbidden (Bukan Admin)"
// @Security     BearerAuth
// @Router       /admin/loans [get]
func GetAllLoans(c *gin.Context) {
	var loans []models.Loan

	if err := database.DB.Preload("User").Preload("Book").
		Order("created_at DESC").
		Find(&loans).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch all loans", err.Error())
		return
	}

	helpers.SuccessResponse(c, "All loan data", loans)
}