package users

import (
	"log"

	"gin-gonic/helper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserServer struct {
	router  *gin.RouterGroup
	db      *gorm.DB
	version string
}

func NewUserServer(router *gin.RouterGroup, db *gorm.DB, version string) *UserServer {
	return &UserServer{router: router, db: db, version: version}
}

func (s *UserServer) Init() {
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	if config.AUTO_MIGRATE == "Y" {
		if err := s.db.AutoMigrate(&User{}); err != nil {
			log.Printf("Failed to auto migrate User: %v", err)
		}
	}

	service := NewUserService(s.db)
	controller := NewUserController(service)

	// Auth routes
	auth := s.router.Group("/" + s.version + "/auth")
	auth.POST("/register", controller.Create)
	auth.POST("/login", controller.Login)

	// Protected user routes
	userRoutes := s.router.Group("/" + s.version + "/users")
	userRoutes.GET("/profile", controller.GetProfile)
	userRoutes.PUT("/:id", controller.Update)

	// Admin user management
	adminUsers := s.router.Group("/" + s.version + "/admin/users")
	adminUsers.GET("/users", controller.GetList)
	adminUsers.GET("/all", controller.GetList2)
	adminUsers.GET("/search", controller.Search)
	adminUsers.GET("/:id", controller.GetByID)
	adminUsers.DELETE("/:id", controller.Delete)
	adminUsers.GET("/stats", controller.GetStats)
}
