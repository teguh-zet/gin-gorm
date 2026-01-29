package loans

type LoanNatsController interface {
	ProcessBorrow(payload PayloadLoan)
	ProcessReturn(payload PayloadLoan)
}

type loanNatsController struct {
	service LoanNatsService
}

func NewLoanNatsController(service LoanNatsService) LoanNatsController {
	return &loanNatsController{service: service}
}

func (c *loanNatsController) ProcessBorrow(payload PayloadLoan) {
	c.service.ProcessBorrow(payload)
}

func (c *loanNatsController) ProcessReturn(payload PayloadLoan) {
	c.service.ProcessReturn(payload)
}
