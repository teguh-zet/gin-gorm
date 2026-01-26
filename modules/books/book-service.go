package books

import (
	"strconv"

	"gin-gonic/helper"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookService interface {
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

type bookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) BookService {
	return &bookService{db: db}
}

func (s *bookService) GetList(c *gin.Context) {
	var book []Book
	if err := s.db.Find(&book).Error; err != nil {
		helper.InternalServerError(c, "Failed to fetch book", err.Error())
		return
	}
	helper.SuccessResponse(c, "Book retrieved successfully", book)
}

func (s *bookService) GetList2(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sortBy := c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "DESC")
	available := c.DefaultQuery("available", "false")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}
	if limitNum > 100 {
		limitNum = 100
	}

	allowedSort := map[string]bool{
		"id":     true,
		"title":  true,
		"author": true,
		"stock":  true,
	}
	if !allowedSort[sortBy] {
		sortBy = "id"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	offset := (pageNum - 1) * limitNum

	query := s.db.Model(&Book{})
	if available == "true" {
		query = query.Where("stock > ?", 0)
	}

	var books []Book
	var total int64

	query.Count(&total)
	if err := query.Offset(offset).Limit(limitNum).Order(sortBy + " " + order).Find(&books).Error; err != nil {
		helper.InternalServerError(c, "Failed to fetch books", err.Error())
		return
	}

	helper.SuccessResponse(c, "Books retrieved successfully", gin.H{
		"data":  books,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

func (s *bookService) GetByID(c *gin.Context) {
	id := c.Param("id")

	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		helper.ErrorResponse(c, 404, "Book not found", err.Error())
		return
	}
	helper.SuccessResponse(c, "Book retrieved successfully", book)
}

func (s *bookService) Create(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationError(c, err.Error())
		return
	}

	book := Book{
		Title:       req.Title,
		Author:      req.Author,
		Stock:       req.Stock,
		BorrowCount: 0,
	}

	if err := s.db.Create(&book).Error; err != nil {
		helper.InternalServerError(c, "Failed to create book", err.Error())
		return
	}

	helper.CreatedResponse(c, "Book successfully created", book)
}

func (s *bookService) Update(c *gin.Context) {
	id := c.Param("id")

	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		helper.ErrorResponse(c, 404, "Book not found", err.Error())
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationError(c, err.Error())
		return
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Stock != 0 {
		book.Stock = req.Stock
	}

	if err := s.db.Save(&book).Error; err != nil {
		helper.InternalServerError(c, "Failed to update book", err.Error())
		return
	}

	helper.SuccessResponse(c, "Book updated successfully", book)
}

func (s *bookService) Delete(c *gin.Context) {
	id := c.Param("id")

	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		helper.ErrorResponse(c, 404, "Book not found", err.Error())
		return
	}

	if err := s.db.Delete(&book).Error; err != nil {
		helper.InternalServerError(c, "Failed to delete book", err.Error())
		return
	}

	helper.SuccessResponse(c, "Book deleted successfully", book)
}

func (s *bookService) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationError(c, err.Error())
		return
	}

	if err := s.db.Delete(&[]Book{}, req.IDs).Error; err != nil {
		helper.InternalServerError(c, "Failed to delete books", err.Error())
		return
	}

	helper.SuccessResponse(c, "Books deleted successfully", req.IDs)
}

func (s *bookService) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		helper.BadRequestError(c, "Query parameter q is required", nil)
		return
	}

	var books []Book
	if err := s.db.Where("title ILIKE ? OR author ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&books).Error; err != nil {
		helper.InternalServerError(c, "Failed to search books", err.Error())
		return
	}

	helper.SuccessResponse(c, "Search results", books)
}

func (s *bookService) UploadImage(c *gin.Context) {
	id := c.Param("id")

	var book Book
	if err := s.db.First(&book, id).Error; err != nil {
		helper.NotFoundError(c, "Book not found")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		helper.BadRequestError(c, "Image file is required", err.Error())
		return
	}

	config, err := helper.LoadConfig(".")
	if err != nil {
		helper.InternalServerError(c, "Failed to load config", err.Error())
		return
	}

	if config.CloudinaryCloudName == "" || config.CloudinaryAPIKey == "" || config.CloudinaryAPISecret == "" {
		helper.InternalServerError(c, "Cloudinary config missing", "Please set CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, CLOUDINARY_API_SECRET")
		return
	}

	cld, err := cloudinary.NewFromParams(
		config.CloudinaryCloudName,
		config.CloudinaryAPIKey,
		config.CloudinaryAPISecret,
	)
	if err != nil {
		helper.InternalServerError(c, "Failed to initialize Cloudinary", err.Error())
		return
	}

	uploadResult, err := cld.Upload.Upload(c, file, uploader.UploadParams{})
	if err != nil {
		helper.InternalServerError(c, "Failed to upload image", err.Error())
		return
	}

	book.ImageURL = uploadResult.SecureURL
	if err := s.db.Save(&book).Error; err != nil {
		helper.InternalServerError(c, "Failed to save image URL", err.Error())
		return
	}

	helper.SuccessResponse(c, "Book image uploaded successfully", book)
}
