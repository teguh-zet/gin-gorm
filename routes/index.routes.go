package routes

import (
	book_controller "gin-gonic/controllers/book_controllers"
	loan_controller "gin-gonic/controllers/loan_controllers"
	"gin-gonic/controllers/user_controllers"
	"gin-gonic/middlewares"

	"github.com/gin-gonic/gin"
	
	// Jangan lupa import swagger jika sudah ada
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "gin-gonic/docs" 
)

func InitRoute(app *gin.Engine) {
	// 0. Swagger Route (Bisa diakses publik)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := app.Group("/api/v1")

	// ==========================================
	// 1. PUBLIC ROUTES (Tidak butuh Login)
	// ==========================================
	{
		// Auth
		auth := api.Group("/auth")
		auth.POST("/register", user_controllers.CreateUser)
		auth.POST("/login", user_controllers.Login)
		
		// Public Book Access (Hanya boleh LIHAT buku)
		booksPublic := api.Group("/books")
		booksPublic.GET("", book_controller.GetAllBooks)       // User biasa hanya boleh lihat
		booksPublic.GET("/all", book_controller.GetAllBooks2)  // Pagination
		booksPublic.GET("/search", book_controller.SearchBooks)
		booksPublic.GET("/:id", book_controller.GetBookByID)
	}

	// ==========================================
	// 2. PROTECTED ROUTES (Butuh Token User/Admin)
	// ==========================================
	protected := api.Group("")
	protected.Use(middlewares.JWTMiddleware())
	{
		// --- User Routes (Profile & Edit diri sendiri) ---
		userRoutes := protected.Group("/users")
		userRoutes.GET("/profile", user_controllers.GetProfile)
		userRoutes.PUT("/:id", user_controllers.UpdateUser) // Sebaiknya validasi di controller agar hanya bisa edit diri sendiri
		
		// --- Loan Routes (Peminjaman) ---
		loanRoutes := protected.Group("/loans")
		loanRoutes.POST("/", loan_controller.BorrowBook)
		loanRoutes.GET("/my", loan_controller.GetMyLoans)
		loanRoutes.POST("/return/:id", loan_controller.ReturnBook)
		loanRoutes.GET("fav",loan_controller.GetPopularBooks)

		// ==========================================
		// 3. ADMIN ROUTES (Hanya Role Admin)
		// ==========================================
		admin := protected.Group("/admin")
		admin.Use(middlewares.AdminMiddleware())
		{
			// Admin User Management
			adminUsers := admin.Group("/users")
			adminUsers.GET("/users", user_controllers.GetAllUsers)
			adminUsers.GET("/all", user_controllers.GetAllUsers2)
			adminUsers.GET("/search", user_controllers.SearchUsers)
			adminUsers.GET("/:id", user_controllers.GetUserByID)
			adminUsers.DELETE("/:id", user_controllers.DeleteUser)
			// admin.DELETE("/users/bulk", user_controllers.BulkDeleteUsers)

			// Admin Book Management (CRUD Buku pindah ke sini agar aman)
			adminBooks := admin.Group("/books")
			adminBooks.POST("", book_controller.CreateBook)
			adminBooks.PUT("/:id", book_controller.UpdateBook)
			adminBooks.DELETE("/:id", book_controller.DeleteBook)
			adminBooks.DELETE("/bulk-delete", book_controller.BulkDeleteBooks)
			adminBooks.GET("/stats", loan_controller.GetLoanStats)
			adminBooks.PATCH("/:id/image", book_controller.UploadBookImage)
			// Admin Loan Monitoring (Optional: jika ingin melihat semua pinjaman)
			// admin.GET("/loans", loan_controller.GetAllLoans) 
		}
	}
}