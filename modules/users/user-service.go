package users

import (
	"net/http"
	"strconv"
	"time"

	"gin-gonic/helper"
	"gin-gonic/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService interface {
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

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) GetList2(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	if pageNum < 1 {
		pageNum = 1
	}
	if limitNum < 1 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum

	var users []User
	var total int64

	if err := s.db.Offset(offset).Limit(limitNum).Find(&users).Error; err != nil {
		helper.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	s.db.Model(&User{}).Count(&total)

	helper.SuccessResponse(c, "Users retrieved successfully", gin.H{
		"data":  users,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

func (s *userService) GetList(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "DESC")

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
		"id":         true,
		"name":       true,
		"email":      true,
		"address":    true,
		"born_date":  true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	offset := (pageNum - 1) * limitNum

	var users []User
	var total int64

	query := s.db.Model(&User{})
	query.Count(&total)

	if err := query.Offset(offset).Limit(limitNum).Order(sortBy + " " + order).Find(&users).Error; err != nil {
		helper.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	helper.SuccessResponse(c, "Users retrieved successfully", gin.H{
		"data":  users,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

func (s *userService) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helper.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.NotFoundError(c, "User not found")
			return
		}
		helper.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	helper.SuccessResponse(c, "User retrieved successfully", user)
}

func (s *userService) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationError(c, err.Error())
		return
	}

	bornDate, err := time.Parse("2006-01-02", req.BornDate)
	if err != nil {
		helper.BadRequestError(c, "Invalid date format", "born_date must be in YYYY-MM-DD format")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		helper.InternalServerError(c, "Failed to process password", err.Error())
		return
	}

	user := User{
		Name:     req.Name,
		Address:  req.Address,
		Email:    req.Email,
		Password: hashedPassword,
		BornDate: bornDate,
		Role:     "user",
	}

	if err := s.db.Create(&user).Error; err != nil {
		helper.InternalServerError(c, "Failed to create user", err.Error())
		return
	}

	helper.CreatedResponse(c, "User created successfully", user)
}

func (s *userService) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helper.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.NotFoundError(c, "User not found")
			return
		}
		helper.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationError(c, err.Error())
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Address != "" {
		user.Address = req.Address
	}
	if req.Email != "" {
		var existing User
		if err := s.db.Where("email = ? AND id != ?", req.Email, id).First(&existing).Error; err == nil {
			helper.ErrorResponse(c, http.StatusConflict, "Email already exists", "Email must be unique")
			return
		}
		user.Email = req.Email
	}
	if req.BornDate != "" {
		bornDate, err := time.Parse("2006-01-02", req.BornDate)
		if err != nil {
			helper.BadRequestError(c, "Invalid date format", "born_date must be in YYYY-MM-DD format")
			return
		}
		user.BornDate = bornDate
	}

	if err := s.db.Save(&user).Error; err != nil {
		helper.InternalServerError(c, "Failed to update user", err.Error())
		return
	}

	helper.SuccessResponse(c, "User updated successfully", user)
}

func (s *userService) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helper.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.NotFoundError(c, "User not found")
			return
		}
		helper.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	if err := s.db.Delete(&user).Error; err != nil {
		helper.InternalServerError(c, "Failed to delete user", err.Error())
		return
	}

	helper.SuccessResponse(c, "User deleted successfully", gin.H{"id": id})
}

func (s *userService) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		helper.BadRequestError(c, "Query parameter q is required", nil)
		return
	}

	var users []User
	if err := s.db.Where("name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&users).Error; err != nil {
		helper.InternalServerError(c, "Failed to search users", err.Error())
		return
	}

	helper.SuccessResponse(c, "Search results", users)
}

func (s *userService) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationError(c, err.Error())
		return
	}

	var user User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		helper.InternalServerError(c, "Failed to generate token", err.Error())
		return
	}

	loginResponse := LoginResponse{Token: token, User: user}
	helper.SuccessResponse(c, "Login successful", loginResponse)
}

func (s *userService) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userIDUint, ok := userID.(float64)
	if !ok {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	var user User
	if err := s.db.First(&user, uint(userIDUint)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			helper.NotFoundError(c, "User not found")
			return
		}
		helper.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	helper.SuccessResponse(c, "Profile retrieved successfully", user)
}

type UserStats struct {
	TotalUsers    int64  `json:"total_users"`
	NewUsersToday int64  `json:"new_users_today"`
	LatestUsers   []User `json:"latest_users"`
}

func (s *userService) GetStats(c *gin.Context) {
	var totalUsers int64
	var newUsersToday int64
	var latestUsers []User

	s.db.Model(&User{}).Count(&totalUsers)

	today := time.Now().Format("2006-01-02")
	s.db.Where("DATE(created_at) = ?", today).Model(&User{}).Count(&newUsersToday)

	s.db.Order("created_at DESC").Limit(5).Find(&latestUsers)

	stats := UserStats{
		TotalUsers:    totalUsers,
		NewUsersToday: newUsersToday,
		LatestUsers:   latestUsers,
	}

	helper.SuccessResponse(c, "User statistics", stats)
}
