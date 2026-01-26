package loans

import "github.com/gin-gonic/gin"

// Controller hanya memanggil service (logika bisnis ada di service)
func GetLoanStats(c *gin.Context)    { GetLoanStatsService(c) }
func GetPopularBooks(c *gin.Context) { GetPopularBooksService(c) }
func BorrowBook(c *gin.Context)      { BorrowBookService(c) }
func ReturnBook(c *gin.Context)      { ReturnBookService(c) }
func GetMyLoans(c *gin.Context)      { GetMyLoansService(c) }
func GetAllLoans(c *gin.Context)     { GetAllLoansService(c) }
