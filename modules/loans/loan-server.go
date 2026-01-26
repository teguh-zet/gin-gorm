package loans

import (
	"log"

	"gin-gonic/helper"
	"gin-gonic/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoanServer struct {
	router  *gin.RouterGroup
	db      *gorm.DB
	version string
}

func NewLoanServer(router *gin.RouterGroup, db *gorm.DB, version string) *LoanServer {
	return &LoanServer{router: router, db: db, version: version}
}

func (s *LoanServer) Init() {
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	if config.AUTO_MIGRATE == "Y" {
		if err := s.db.AutoMigrate(&Loan{}); err != nil {
			log.Printf("Failed to auto migrate Loan: %v", err)
		}
	}

	service := NewLoanService(s.db)
	controller := NewLoanController(service)

	// Protected loan routes
	loanRoutes := s.router.Group("/" + s.version + "/loans")
	loanRoutes.Use(middlewares.JWTMiddleware())
	loanRoutes.POST("/", controller.Borrow)
	loanRoutes.GET("/my", controller.GetMy)
	loanRoutes.POST("/return/:id", controller.Return)
	loanRoutes.GET("/fav", controller.GetPopularBooks)

	// Admin loan stats
	adminRoutes := s.router.Group("/" + s.version + "/admin")
	adminRoutes.Use(middlewares.JWTMiddleware(), middlewares.AdminMiddleware())
	adminRoutes.GET("/books/stats", controller.GetStats)
	adminRoutes.GET("/loans", controller.GetAll)
}
