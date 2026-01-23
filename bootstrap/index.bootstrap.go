package bootstrap

import (
	"fmt"
	"gin-gonic/database"
	"gin-gonic/helpers"
	"gin-gonic/models"
	"gin-gonic/routes"
	"gin-gonic/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func BootstrapApp() {
	database.ConnectDatabase()

	// Smart auto migration dengan handling untuk existing data
	fmt.Println("Running smart auto migration...")
	err := runSmartMigration()
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	fmt.Println("Auto migration completed successfully")

	app := gin.Default()
	app.Static("/uploads","./uploads")
	seedAdmin()
	// Middleware untuk logging dan recovery
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	routes.InitRoute(app)

	fmt.Printf("Server starting on port %s\n", helpers.GetConfig("APP_PORT"))
	app.Run(helpers.GetConfig("APP_PORT"))
}

// runSmartMigration menangani migration dengan aman untuk existing data
func runSmartMigration() error {
	// Cek apakah tabel users sudah ada dan memiliki data
	var userCount int64
	hasUsersTable := database.DB.Migrator().HasTable(&models.User{})
	if hasUsersTable {
		database.DB.Model(&models.User{}).Count(&userCount)
	}

	// Jika ada users existing tanpa password, beri password default
	if hasUsersTable && userCount > 0 {
		fmt.Printf("Found %d existing users, checking for password migration...\n", userCount)

		// Cek apakah kolom password sudah ada
		hasPasswordColumn := database.DB.Migrator().HasColumn(&models.User{}, "password")

		if !hasPasswordColumn {
			fmt.Println("Adding password column to existing users table...")
			// GORM akan menangani penambahan kolom nullable terlebih dahulu
		}

		// Cari users yang belum punya password
		var usersWithoutPassword []models.User
		database.DB.Where("password = '' OR password IS NULL").Find(&usersWithoutPassword)

		if len(usersWithoutPassword) > 0 {
			fmt.Printf("Found %d users without passwords, setting default password...\n", len(usersWithoutPassword))

			defaultPassword := "defaultpassword123"
			hashedPassword, err := utils.HashPassword(defaultPassword)
			if err != nil {
				return fmt.Errorf("failed to hash default password: %v", err)
			}

			// Update password untuk users yang belum punya
			result := database.DB.Model(&models.User{}).
				Where("password = '' OR password IS NULL").
				Update("password", hashedPassword)

			if result.Error != nil {
				return fmt.Errorf("failed to update user passwords: %v", result.Error)
			}

			fmt.Printf(" Successfully migrated %d users with default password\n", result.RowsAffected)
			fmt.Printf(" Default password: %s\n", defaultPassword)
			fmt.Println(" IMPORTANT: Users should change their password after first login!")
		} else {
			fmt.Println("All existing users already have passwords")
		}
	}

	// Jalankan auto migration normal
	fmt.Println("Running GORM auto migration...")
	err := database.DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Loan{})
	if err != nil {
		return fmt.Errorf("auto migration failed: %v", err)
	}

	fmt.Println("Database migration completed successfully")
	return nil
}


func seedAdmin(){
	adminEmail := helpers.GetConfig("ADMIN_EMAIL")
	adminPassword := helpers.GetConfig("ADMIN_PASSWORD")
	adminName := helpers.GetConfig("ADMIN_NAME")

	// validasi env tidak diset
	if adminEmail =="" ||adminPassword ==""{
		fmt.Println("seeding skipped: ADMIN_EMAIL or ADMIN_PASSWORD not found in .env")
	}
	var count int64

	//cek apakah admin sudah ada
	database.DB.Model(&models.User{}).Where("role =?", "admin").Count(&count)

	if count == 0{
		fmt.Println("No admin found. Creating admin from environment variable")
		//hash password dari env
		hashedPassword, err := utils.HashPassword(adminPassword)
		if err != nil{
			fmt.Printf("Error hashing password : %v\n",err)
			return
		}
		admin := models.User{
			Name:      adminName,
			Email:     adminEmail,
			Password:  hashedPassword,
			Address:   "System Administrator",
			Role:      "admin", 
			BornDate:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := database.DB.Create(&admin).Error;err!=nil{
			fmt.Printf("failed to create admin :%v \n",err)
		}else{
			fmt.Println("admin account seeded succesfully")
		}

	
	
	}else{
		fmt.Println("Admin account check: OK (Admin already exists")
	}

}