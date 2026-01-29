package books

import (
	"log"

	"gin-gonic/helper"
	"gin-gonic/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookServer struct {
	router  *gin.RouterGroup
	db      *gorm.DB
	version string
}

func NewBookServer(router *gin.RouterGroup, db *gorm.DB, version string) *BookServer {
	return &BookServer{router: router, db: db, version: version}
}

func (s *BookServer) Init() {
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	if config.AUTO_MIGRATE == "Y" {
		if err := s.db.AutoMigrate(&Book{}); err != nil {
			log.Printf("Failed to auto migrate Book: %v", err)
		}
	}

	service := NewBookService(s.db)
	controller := NewBookController(service)

	// Public book routes
	booksPublic := s.router.Group("/" + s.version + "/books")
	booksPublic.GET("", controller.GetList)
	booksPublic.GET("/all", controller.GetList2)
	booksPublic.GET("/search", controller.Search)
	booksPublic.GET("/:id", controller.GetByID)

	// Admin book routes (protected)
	adminBooks := s.router.Group("/" + s.version + "/admin/books")
	adminBooks.Use(middlewares.JWTMiddleware(), middlewares.AdminMiddleware())
	adminBooks.POST("", controller.Create)
	adminBooks.PUT("/:id", controller.Update)
	adminBooks.DELETE("/:id", controller.Delete)
	adminBooks.DELETE("/bulk-delete", controller.BulkDelete)
	adminBooks.PATCH("/:id/image", controller.UploadImage)
}
