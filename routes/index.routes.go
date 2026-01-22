package routes

import (
	book_controller "gin-gonic/controllers/book_controllers"
	"gin-gonic/controllers/user_controllers"
	"gin-gonic/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoute(app *gin.Engine) {
	// API versioning
	api := app.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", user_controllers.Login) // POST /api/v1/auth/login
		}

		// Protected routes (require JWT token)
		protected := api.Group("")
		protected.Use(middlewares.JWTMiddleware())
		{
			// Protected user routes
			protectedUsers := protected.Group("/users")
			{
				protectedUsers.GET("/stats", user_controllers.GetUserStats)      // GET /api/v1/users/stats
				protectedUsers.GET("/profile", user_controllers.GetProfile) // GET /api/v1/users/profile
			}
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("", user_controllers.GetAllUsers)             // GET /api/v1/users
			users.GET("/all", user_controllers.GetAllUsers2)        // GET /api/v1/users/all <-- untuk pagination
			users.GET("/:id", user_controllers.GetUserByID)         // GET /api/v1/users/:id
			users.POST("", user_controllers.CreateUser)             // POST /api/v1/users
			users.PUT("/:id", user_controllers.UpdateUser)          // PUT /api/v1/users/:id
			users.DELETE("/:id", user_controllers.DeleteUser)       // DELETE /api/v1/users/:id
			users.GET("/search", user_controllers.SearchUsers)      // GET /api/v1/users/search?q=john
			users.DELETE("/bulk", user_controllers.BulkDeleteUsers) // DELETE /api/v1/users/bulk
		}

		// Book routes (placeholder)
		books := api.Group("/books")
		{
			books.GET("/search", book_controller.SearchBooks)
			books.DELETE("/bulk-delete", book_controller.BulkDeleteBooks)
			books.GET("/all", book_controller.GetAllBooks2)
			books.GET("/:id", book_controller.GetBookByID)
			books.DELETE("/:id", book_controller.DeleteBook)
			books.PUT("/:id", book_controller.UpdateBook)
			books.POST("", book_controller.CreateBook)
			books.GET("", book_controller.GetAllBooks)
		}
	}

}
