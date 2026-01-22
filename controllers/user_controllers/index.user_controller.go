package user_controllers

import (
	"net/http"
	"strconv"
	"time"

	"gin-gonic/database"
	"gin-gonic/helpers"
	"gin-gonic/models"
	"gin-gonic/utils"

	"github.com/gin-gonic/gin"
)

// GetAllUsers mengambil semua data user dengan pagination
func GetAllUsers2(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)

	offset := (pageNum - 1) * limitNum

	var users []models.User
	var total int64

	if err := database.DB.Offset(offset).Limit(limitNum).Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	database.DB.Model(&models.User{}).Count(&total)

	helpers.SuccessResponse(c, "Users retrieved successfully", gin.H{
		"data":  users,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

// with pagination, sorting and with a little validation
func GetAllUsers(c *gin.Context) {
	// Pagination parameters
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Sorting parameters
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "DESC")

	// Validasi page dan limit
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

	offset := (pageNum - 1) * limitNum

	// Validasi sortBy
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

	var users []models.User
	var total int64

	if err := database.DB.
		Order(sortBy + " " + order).
		Offset(offset).
		Limit(limitNum).
		Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	database.DB.Model(&models.User{}).Count(&total)

	totalPages := (total + int64(limitNum) - 1) / int64(limitNum)

	helpers.SuccessResponse(c, "Users retrieved successfully", gin.H{
		"data": users,
		"pagination": gin.H{
			"total":        total,
			"page":         pageNum,
			"limit":        limitNum,
			"total_pages":  totalPages,
			"has_next":     pageNum < int(totalPages),
			"has_previous": pageNum > 1,
		},
		"sorting": gin.H{
			"sort_by": sortBy,
			"order":   order,
		},
	})
}

// with sorting
func GetAllUsers3(c *gin.Context) {
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "ASC") // ASC atau DESC

	var users []models.User

	if err := database.DB.Order(sortBy + " " + order).Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	helpers.SuccessResponse(c, "Users retrieved successfully", users)
}

// GetUserByID mengambil user berdasarkan ID
func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "User not found")
			return
		}
		helpers.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	helpers.SuccessResponse(c, "User retrieved successfully", user)
}

// CreateUser membuat user baru
func CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err.Error())
		return
	}

	// Parse tanggal lahir
	bornDate, err := time.Parse("2006-01-02", req.BornDate)
	if err != nil {
		helpers.BadRequestError(c, "Invalid date format", "born_date must be in YYYY-MM-DD format")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		helpers.InternalServerError(c, "Failed to process password", err.Error())
		return
	}

	user := models.User{
		Name:     req.Name,
		Address:  req.Address,
		Email:    req.Email,
		Password: hashedPassword,
		BornDate: bornDate,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		// Check for duplicate email
		if err.Error() == "UNIQUE constraint failed: users.email" || err.Error() == "duplicate key value violates unique constraint \"users_email_key\"" {
			helpers.ErrorResponse(c, http.StatusConflict, "Email already exists", "Email must be unique")
			return
		}
		helpers.InternalServerError(c, "Failed to create user", err.Error())
		return
	}

	helpers.CreatedResponse(c, "User created successfully", user)
}

// UpdateUser mengupdate data user
func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "User not found")
			return
		}
		helpers.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err.Error())
		return
	}

	// Update hanya field yang dikirim
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Email != "" {
		// Check if email already exists for other users
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			helpers.ErrorResponse(c, http.StatusConflict, "Email already exists", "Email must be unique")
			return
		}
		updates["email"] = req.Email
	}
	if req.BornDate != "" {
		bornDate, err := time.Parse("2006-01-02", req.BornDate)
		if err != nil {
			helpers.BadRequestError(c, "Invalid date format", "born_date must be in YYYY-MM-DD format")
			return
		}
		updates["born_date"] = bornDate
	}

	if len(updates) == 0 {
		helpers.BadRequestError(c, "No fields to update", "At least one field must be provided")
		return
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		helpers.InternalServerError(c, "Failed to update user", err.Error())
		return
	}

	// Fetch updated user
	database.DB.First(&user, id)
	helpers.SuccessResponse(c, "User updated successfully", user)
}

// DeleteUser menghapus user (soft delete)
func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "User not found")
			return
		}
		helpers.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		helpers.InternalServerError(c, "Failed to delete user", err.Error())
		return
	}

	helpers.SuccessResponse(c, "User deleted successfully", gin.H{"id": id})
}

// bulk delete
func BulkDeleteUsers(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		helpers.ValidationError(c, err.Error())
		return
	}
	if err := database.DB.Delete(&[]models.User{}, req.IDs).Error; err != nil {
		helpers.InternalServerError(c, " failed to delete users", err.Error())
		return
	}
	helpers.SuccessResponse(c, "users deleted Succesfully",
		gin.H{"delete_count": len(req.IDs)})

}

// search

func SearchUsers(c *gin.Context) {
	query := c.Query("name")
	if query == "" {
		helpers.BadRequestError(c, "search query required", "Parameter 'name' harus diisi")
		return
	}
	var users []models.User
	if err := database.DB.Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "failed to search users", err.Error())
		return
	}
	helpers.SuccessResponse(c, "search results", users)

}

//user staticstic

type UserStats struct {
	TotalUsers    int64         `json:"total_users"`
	NewUsersToday int64         `json:"new_users_today"`
	LatestUsers   []models.User `json:"latest_users"`
}

func GetUserStats(c *gin.Context) {
	var totalUsers int64
	var newUsersToday int64
	var latestUsers []models.User

	// Total users
	database.DB.Model(&models.User{}).Count(&totalUsers)

	// Users created today
	today := time.Now().Format("2006-01-02")
	database.DB.Where("DATE(created_at) = ?", today).
		Model(&models.User{}).Count(&newUsersToday)

	// Latest 5 users
	database.DB.Order("created_at DESC").Limit(5).Find(&latestUsers)

	stats := UserStats{
		TotalUsers:    totalUsers,
		NewUsersToday: newUsersToday,
		LatestUsers:   latestUsers,
	}
	helpers.SuccessResponse(c, "User statistics", stats)
}

// Login melakukan autentikasi user dan mengembalikan JWT token
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err.Error())
		return
	}

	// Cari user berdasarkan email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
			return
		}
		helpers.InternalServerError(c, "Failed to authenticate", err.Error())
		return
	}

	// Verifikasi password
	if !utils.CheckPassword(req.Password, user.Password) {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		helpers.InternalServerError(c, "Failed to generate token", err.Error())
		return
	}

	// Response dengan token dan user info
	loginResponse := models.LoginResponse{
		Token: token,
		User:  user,
	}

	helpers.SuccessResponse(c, "Login successful", loginResponse)
}

// GetProfile mengambil profile user yang sedang login (butuh JWT token)
func GetProfile(c *gin.Context) {
	// Ambil user ID dari JWT token (sudah di-set oleh middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		helpers.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userIDUint, ok := userID.(float64)
	if !ok {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Invalid user ID", nil)
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(userIDUint)).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "User not found")
			return
		}
		helpers.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	helpers.SuccessResponse(c, "Profile retrieved successfully", user)
}

//
