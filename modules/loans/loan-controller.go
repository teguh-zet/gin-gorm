package loans

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoanController interface {
	GetStats(ctx *gin.Context)
	GetPopularBooks(ctx *gin.Context)
	Borrow(ctx *gin.Context)
	Return(ctx *gin.Context)
	GetMy(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}

type loanController struct {
	service LoanService
}

func NewLoanController(service LoanService) LoanController {
	return &loanController{service: service}
}

func (c *loanController) GetStats(ctx *gin.Context) {
	stats, err := c.service.GetStats()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, stats)
}

func (c *loanController) GetPopularBooks(ctx *gin.Context) {
	books, err := c.service.GetPopularBooks()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": books})
}

func (c *loanController) Borrow(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDFloat, ok := userIDVal.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var input LoanRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	loan, err := c.service.Borrow(uint(userIDFloat), &input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loan)
}

func (c *loanController) Return(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Return(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}

func (c *loanController) GetMy(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDFloat, ok := userIDVal.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	loans, err := c.service.GetMy(uint(userIDFloat))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": loans})
}

func (c *loanController) GetAll(ctx *gin.Context) {
	loans, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": loans})
}
