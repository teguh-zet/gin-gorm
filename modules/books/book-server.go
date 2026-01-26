package books

import "github.com/gin-gonic/gin"

// RegisterRoutes memasang semua route untuk modul buku
func RegisterRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	// Public book access
	booksPublic := api.Group("/books")
	booksPublic.GET("", GetAllBooks)
	booksPublic.GET("/all", GetAllBooks2)
	booksPublic.GET("/search", SearchBooks)
	booksPublic.GET("/:id", GetBookByID)

	// Admin book management
	adminBooks := admin.Group("/books")
	adminBooks.POST("", CreateBook)
	adminBooks.PUT("/:id", UpdateBook)
	adminBooks.DELETE("/:id", DeleteBook)
	adminBooks.DELETE("/bulk-delete", BulkDeleteBooks)
	adminBooks.PATCH("/:id/image", UploadBookImage)
}
