package users

import "github.com/gin-gonic/gin"

// RegisterRoutes memasang semua route untuk modul user
func RegisterRoutes(api *gin.RouterGroup, protected *gin.RouterGroup, admin *gin.RouterGroup) {
	// Auth
	auth := api.Group("/auth")
	auth.POST("/register", CreateUser)
	auth.POST("/login", Login)

	// User protected routes
	userRoutes := protected.Group("/users")
	userRoutes.GET("/profile", GetProfile)
	userRoutes.PUT("/:id", UpdateUser)

	// Admin user management
	adminUsers := admin.Group("/users")
	adminUsers.GET("/users", GetAllUsers)
	adminUsers.GET("/all", GetAllUsers2)
	adminUsers.GET("/search", SearchUsers)
	adminUsers.GET("/:id", GetUserByID)
	adminUsers.DELETE("/:id", DeleteUser)
}
