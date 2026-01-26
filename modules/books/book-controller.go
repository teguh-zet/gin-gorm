package books

import "github.com/gin-gonic/gin"

// Controller hanya memanggil service (logika bisnis ada di service)
func GetAllBooks(c *gin.Context)     { GetAllBooksService(c) }
func GetAllBooks2(c *gin.Context)    { GetAllBooks2Service(c) }
func GetBookByID(c *gin.Context)     { GetBookByIDService(c) }
func CreateBook(c *gin.Context)      { CreateBookService(c) }
func UpdateBook(c *gin.Context)      { UpdateBookService(c) }
func DeleteBook(c *gin.Context)      { DeleteBookService(c) }
func BulkDeleteBooks(c *gin.Context) { BulkDeleteBooksService(c) }
func SearchBooks(c *gin.Context)     { SearchBooksService(c) }
func UploadBookImage(c *gin.Context) { UploadBookImageService(c) }
