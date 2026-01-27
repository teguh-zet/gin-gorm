package books

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookController interface {
	GetList(ctx *gin.Context)
	GetList2(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BulkDelete(ctx *gin.Context)
	Search(ctx *gin.Context)
	UploadImage(ctx *gin.Context)
}

type bookController struct {
	service BookService
}

func NewBookController(service BookService) BookController {
	return &bookController{service: service}
}

func (c *bookController) GetList(ctx *gin.Context) {
	books, err := c.service.GetList()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, books)
}

func (c *bookController) GetList2(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	sortBy := ctx.DefaultQuery("sort_by", "id")
	order := ctx.DefaultQuery("order", "DESC")
	available := ctx.DefaultQuery("available", "false")

	books, total, err := c.service.GetList2(page, limit, sortBy, order, available == "true")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":      books,
		"total_row": total,
	})
}

func (c *bookController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	book, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (c *bookController) Create(ctx *gin.Context) {
	var input CreateBookRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	book, err := c.service.Create(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func (c *bookController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var input UpdateBookRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	book, err := c.service.Update(id, &input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Gagal memperbarui data: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, book)
}

func (c *bookController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menghapus data: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func (c *bookController) BulkDelete(ctx *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	if err := c.service.BulkDelete(req.IDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func (c *bookController) Search(ctx *gin.Context) {
	q := ctx.Query("q")
	books, err := c.service.Search(q)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": books})
}

func (c *bookController) UploadImage(ctx *gin.Context) {
	id := ctx.Param("id")
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required: " + err.Error()})
		return
	}

	book, err := c.service.UploadImage(id, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, book)
}
