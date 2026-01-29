package users

import (
	"errors"
	"fmt"
	"time"

	"gin-gonic/utils"

	"gorm.io/gorm"
)

type UserService interface {
	Create(input *CreateUserRequest) (*User, error)
	Update(id string, input *UpdateUserRequest) (*User, error)
	Delete(id string) error
	GetList(page, limit int, sortBy, order string) ([]User, int64, error)
	GetList2(page, limit int) ([]User, int64, error)
	GetByID(id string) (*User, error)
	Search(query string) ([]User, error)
	Login(input *LoginRequest) (*LoginResponse, error)
	GetProfile(userID uint) (*User, error)
	GetStats() (*UserStats, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

func (s *userService) Create(input *CreateUserRequest) (*User, error) {
	bornDate, err := time.Parse("2006-01-02", input.BornDate)
	if err != nil {
		return nil, errors.New("born_date must be in YYYY-MM-DD format")
	}

	var existing User
	if err := s.db.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:      input.Name,
		Address:   input.Address,
		Email:     input.Email,
		Password:  hashedPassword,
		BornDate:  bornDate,
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(id string, input *UpdateUserRequest) (*User, error) {
	var user User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("data tidak ditemukan")
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Address != "" {
		user.Address = input.Address
	}
	if input.Email != "" {
		var existing User
		if err := s.db.Where("email = ? AND id <> ?", input.Email, id).First(&existing).Error; err == nil {
			return nil, errors.New("email already exists")
		}
		user.Email = input.Email
	}
	if input.BornDate != "" {
		bornDate, err := time.Parse("2006-01-02", input.BornDate)
		if err != nil {
			return nil, errors.New("born_date must be in YYYY-MM-DD format")
		}
		user.BornDate = bornDate
	}

	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) Delete(id string) error {
	var user User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return errors.New("data tidak ditemukan")
	}

	return s.db.Delete(&user).Error
}

func (s *userService) GetList(page, limit int, sortBy, order string) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
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

	offset := (page - 1) * limit

	var users []User
	var total int64

	query := s.db.Model(&User{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order(sortBy + " " + order).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *userService) GetList2(page, limit int) ([]User, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []User
	var total int64

	if err := s.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	s.db.Model(&User{}).Count(&total)

	return users, total, nil
}

func (s *userService) GetByID(id string) (*User, error) {
	var user User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("data tidak ditemukan")
	}

	return &user, nil
}

func (s *userService) Search(query string) ([]User, error) {
	if query == "" {
		return nil, errors.New("query parameter q is required")
	}

	var users []User
	if err := s.db.Where("name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) Login(input *LoginRequest) (*LoginResponse, error) {
	var user User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Name, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token, User: user}, nil
}

func (s *userService) GetProfile(userID uint) (*User, error) {
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("data tidak ditemukan")
	}

	return &user, nil
}

type UserStats struct {
	TotalUsers    int64  `json:"total_users"`
	NewUsersToday int64  `json:"new_users_today"`
	LatestUsers   []User `json:"latest_users"`
}

func (s *userService) GetStats() (*UserStats, error) {
	var totalUsers int64
	var newUsersToday int64
	var latestUsers []User

	if err := s.db.Model(&User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	if err := s.db.Where("DATE(created_at) = ?", today).Model(&User{}).Count(&newUsersToday).Error; err != nil {
		return nil, err
	}

	if err := s.db.Order("created_at DESC").Limit(5).Find(&latestUsers).Error; err != nil {
		return nil, err
	}

	return &UserStats{
		TotalUsers:    totalUsers,
		NewUsersToday: newUsersToday,
		LatestUsers:   latestUsers,
	}, nil
}

func (s *userService) String() string {
	return fmt.Sprintf("userService{db:%v}", s.db)
}
