package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetList(ctx *gin.Context)
	GetList2(ctx *gin.Context)
	GetByID(ctx *gin.Context)
	Search(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	GetStats(ctx *gin.Context)
}

type userController struct {
	service UserService
}

func NewUserController(service UserService) UserController {
	return &userController{service: service}
}

func (c *userController) Create(ctx *gin.Context) {
	var input CreateUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	user, err := c.service.Create(&input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) Update(ctx *gin.Context) {
	var input UpdateUserRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	id := ctx.Param("id")
	user, err := c.service.Update(id, &input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Gagal memperbarui data: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Gagal menghapus data: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Data berhasil dihapus"})
}

func (c *userController) GetList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	order := ctx.DefaultQuery("order", "DESC")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := c.service.GetList(page, limit, sortBy, order)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total_row": total,
	})
}

func (c *userController) GetList2(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := c.service.GetList2(page, limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total_row": total,
	})
}

func (c *userController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak ditemukan: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) Search(ctx *gin.Context) {
	q := ctx.Query("q")
	users, err := c.service.Search(q)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func (c *userController) Login(ctx *gin.Context) {
	var input LoginRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	result, err := c.service.Login(&input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *userController) GetProfile(ctx *gin.Context) {
	userID, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idFloat, ok := userID.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := c.service.GetProfile(uint(idFloat))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *userController) GetStats(ctx *gin.Context) {
	stats, err := c.service.GetStats()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stats)
}
