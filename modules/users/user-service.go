package users

import (
	"net/http"
	"strconv"
	"time"

	"gin-gonic/helpers"
	"gin-gonic/utils"

	"github.com/gin-gonic/gin"
)

// GetAllUsers mengambil semua data user dengan pagination
func GetAllUsers2Service(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)

	offset := (pageNum - 1) * limitNum

	var users []User
	var total int64

	if err := helpers.DB.Offset(offset).Limit(limitNum).Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	helpers.DB.Model(&User{}).Count(&total)

	helpers.SuccessResponse(c, "Users retrieved successfully", gin.H{
		"data":  users,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

// with pagination, sorting and with a little validation
// GetAllUsers godoc
// @Summary      Lihat Semua User (Pagination & Sorting)
// @Description  Menampilkan daftar user dengan fitur pagination, limit, dan sorting. Khusus Admin.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        page      query int    false "Halaman ke berapa (Default: 1)"
// @Param        limit     query int    false "Jumlah data per halaman (Default: 10, Max: 100)"
// @Param        sort_by   query string false "Kolom untuk sorting (id, name, email, created_at, dll). Default: created_at"
// @Param        order     query string false "Arah urutan (ASC atau DESC). Default: DESC"
// @Success      200       {object} map[string]interface{} "Response berisi data array user, objek pagination, dan sorting"
// @Failure      401       {object} map[string]interface{} "Unauthorized"
// @Failure      403       {object} map[string]interface{} "Forbidden (Bukan Admin)"
// @Failure      500       {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /admin/users [get]
func GetAllUsersService(c *gin.Context) {
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

	var users []User
	var total int64

	if err := helpers.DB.
		Order(sortBy + " " + order).
		Offset(offset).
		Limit(limitNum).
		Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	helpers.DB.Model(&User{}).Count(&total)

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
func GetAllUsers3Service(c *gin.Context) {
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "ASC") // ASC atau DESC

	var users []User

	if err := helpers.DB.Order(sortBy + " " + order).Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "Failed to fetch users", err.Error())
		return
	}

	helpers.SuccessResponse(c, "Users retrieved successfully", users)
}

// GetUserByID mengambil user berdasarkan ID
func GetUserByIDService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user User
	if err := helpers.DB.First(&user, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "User not found")
			return
		}
		helpers.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	helpers.SuccessResponse(c, "User retrieved successfully", user)
}

// CreateUser godoc
// @Summary      Register User Baru
// @Description  Mendaftarkan user baru dengan nama, email, password, dan tanggal lahir.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequest true "Data Registrasi User"
// @Success      201  {object} User
// @Failure      400  {object} map[string]interface{} "Validasi Error / Format Tanggal Salah"
// @Failure      409  {object} map[string]interface{} "Email sudah terdaftar"
// @Failure      500  {object} map[string]interface{} "Internal Server Error"
// @Router       /auth/register [post]
func CreateUserService(c *gin.Context) {
	var req CreateUserRequest
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

	user := User{
		Name:     req.Name,
		Address:  req.Address,
		Email:    req.Email,
		Password: hashedPassword,
		BornDate: bornDate,
	}

	if err := helpers.DB.Create(&user).Error; err != nil {
		// Check for duplicate email
		if err.Error() == "UNIQUE constraint failed: email" || err.Error() == "duplicate key value violates unique constraint \"users_email_key\"" {
			helpers.ErrorResponse(c, http.StatusConflict, "Email already exists", "Email must be unique")
			return
		}
		helpers.InternalServerError(c, "Failed to create user", err.Error())
		return
	}

	helpers.CreatedResponse(c, "User created successfully", user)
}

// UpdateUser mengupdate data user
// UpdateUser godoc
// @Summary      Update Data User
// @Description  Mengubah data profil user (Nama, Alamat, Email, Tgl Lahir). Membutuhkan token login.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        id path int true "User ID"
// @Param        request body UpdateUserRequest true "Data yang ingin diupdate"
// @Success      200  {object} User
// @Failure      400  {object} map[string]interface{} "ID Salah / Format Tanggal Salah / Tidak ada field update"
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Failure      404  {object} map[string]interface{} "User tidak ditemukan"
// @Failure      409  {object} map[string]interface{} "Email konflik dengan user lain"
// @Failure      500  {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /users/{id} [put]
func UpdateUserService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user User
	if err := helpers.DB.First(&user, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "User not found")
			return
		}
		helpers.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	var req UpdateUserRequest
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
		var existingUser User
		if err := helpers.DB.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
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

	if err := helpers.DB.Model(&user).Updates(updates).Error; err != nil {
		helpers.InternalServerError(c, "Failed to update user", err.Error())
		return
	}

	// Fetch updated user
	helpers.DB.First(&user, id)
	helpers.SuccessResponse(c, "User updated successfully", user)
}

// DeleteUser menghapus user (soft delete)
// DeleteUser godoc
// @Summary      Hapus User (Admin)
// @Description  Menghapus data user dari database (Soft Delete). Hanya Admin yang boleh melakukan ini.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        id   path int true "User ID yang akan dihapus"
// @Success      200  {object} map[string]interface{} "Mengembalikan ID user yang dihapus"
// @Failure      400  {object} map[string]interface{} "ID Salah / Bukan Angka"
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Failure      403  {object} map[string]interface{} "Forbidden (Bukan Admin)"
// @Failure      404  {object} map[string]interface{} "User tidak ditemukan"
// @Failure      500  {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /admin/users/{id} [delete]
func DeleteUserService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helpers.BadRequestError(c, "Invalid user ID", "ID must be a number")
		return
	}

	var user User
	if err := helpers.DB.First(&user, id).Error; err != nil {
		if err.Error() == "record not found" {
			helpers.NotFoundError(c, "User not found")
			return
		}
		helpers.InternalServerError(c, "Failed to fetch user", err.Error())
		return
	}

	if err := helpers.DB.Delete(&user).Error; err != nil {
		helpers.InternalServerError(c, "Failed to delete user", err.Error())
		return
	}

	helpers.SuccessResponse(c, "User deleted successfully", gin.H{"id": id})
}

// SearchUsers godoc
// @Summary      Mencari User
// @Description  Mencari data user berdasarkan keyword yang cocok dengan nama ATAU email (LIKE query).
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Param        name query string true "Keyword pencarian (Nama atau Email)"
// @Success      200  {array} User
// @Failure      400  {object} map[string]interface{} "Parameter query kosong"
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Failure      500  {object} map[string]interface{} "Internal Server Error"
// @Security     BearerAuth
// @Router       /admin/users/search [get]
func SearchUsersService(c *gin.Context) {
	query := c.Query("name")
	if query == "" {
		helpers.BadRequestError(c, "search query required", "Parameter 'name' harus diisi")
		return
	}
	var users []User
	if err := helpers.DB.Where("name LIKE ? OR email LIKE ?", "%"+query+"%", "%"+query+"%").Find(&users).Error; err != nil {
		helpers.InternalServerError(c, "failed to search users", err.Error())
		return
	}
	helpers.SuccessResponse(c, "search results", users)

}

//user staticstic

type UserStats struct {
	TotalUsers    int64  `json:"total_users"`
	NewUsersToday int64  `json:"new_users_today"`
	LatestUsers   []User `json:"latest_users"`
}

// GetUserStats godoc
// @Summary      Statistik User (Admin)
// @Description  Menampilkan total user, user baru hari ini, dan 5 user terakhir.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200  {object} map[string]interface{} "User Statistics (Total, New Today, Latest)"
// @Failure      401  {object} map[string]interface{} "Unauthorized"
// @Failure      403  {object} map[string]interface{} "Forbidden (Bukan Admin)"
// @Security     BearerAuth
// @Router       /admin/stats [get]
func GetUserStatsService(c *gin.Context) {
	var totalUsers int64
	var newUsersToday int64
	var latestUsers []User

	// Total users
	helpers.DB.Model(&User{}).Count(&totalUsers)

	// Users created today
	today := time.Now().Format("2006-01-02")
	helpers.DB.Where("DATE(created_at) = ?", today).
		Model(&User{}).Count(&newUsersToday)

	// Latest 5 users
	helpers.DB.Order("created_at DESC").Limit(5).Find(&latestUsers)

	stats := UserStats{
		TotalUsers:    totalUsers,
		NewUsersToday: newUsersToday,
		LatestUsers:   latestUsers,
	}
	helpers.SuccessResponse(c, "User statistics", stats)
}

// Login melakukan autentikasi user dan mengembalikan JWT token
// Login godoc
// @Summary      Login User
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Credentials"
// @Success      200  {object} LoginResponse
// @Failure      400  {object} map[string]interface{}
// @Failure      401  {object} map[string]interface{}
// @Router       /auth/login [post]
func LoginService(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationError(c, err.Error())
		return
	}

	// Cari user berdasarkan email
	var user User
	if err := helpers.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
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
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		helpers.InternalServerError(c, "Failed to generate token", err.Error())
		return
	}

	// Response dengan token dan user info
	loginResponse := LoginResponse{
		Token: token,
		User:  user,
	}

	helpers.SuccessResponse(c, "Login successful", loginResponse)
}

// GetProfile mengambil profile user yang sedang login (butuh JWT token)
// GetProfile godoc
// @Summary      Profil Saya
// @Description  Menampilkan data detail user yang sedang login berdasarkan Token.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer Token"
// @Success      200  {object} User
// @Failure      401  {object} map[string]interface{} "Token tidak ada / tidak valid"
// @Failure      404  {object} map[string]interface{} "User tidak ditemukan"
// @Failure      500  {object} map[string]interface{} "Gagal mengambil data"
// @Security     BearerAuth
// @Router       /users/profile [get]
func GetProfileService(c *gin.Context) {
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

	var user User
	if err := helpers.DB.First(&user, uint(userIDUint)).Error; err != nil {
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
