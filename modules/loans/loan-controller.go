package loans

import "github.com/gin-gonic/gin"

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

func (c *loanController) GetStats(ctx *gin.Context)        { c.service.GetStats(ctx) }
func (c *loanController) GetPopularBooks(ctx *gin.Context) { c.service.GetPopularBooks(ctx) }
func (c *loanController) Borrow(ctx *gin.Context)          { c.service.Borrow(ctx) }
func (c *loanController) Return(ctx *gin.Context)          { c.service.Return(ctx) }
func (c *loanController) GetMy(ctx *gin.Context)           { c.service.GetMy(ctx) }
func (c *loanController) GetAll(ctx *gin.Context)          { c.service.GetAll(ctx) }
