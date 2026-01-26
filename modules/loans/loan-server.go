package loans

import "github.com/gin-gonic/gin"

// RegisterRoutes memasang semua route untuk modul loan
func RegisterRoutes(protected *gin.RouterGroup, admin *gin.RouterGroup) {
	// Loan routes (protected)
	loanRoutes := protected.Group("/loans")
	loanRoutes.POST("/", BorrowBook)
	loanRoutes.GET("/my", GetMyLoans)
	loanRoutes.POST("/return/:id", ReturnBook)
	loanRoutes.GET("/fav", GetPopularBooks)

	// Admin loan stats
	admin.GET("/books/stats", GetLoanStats)
}
