package main

import (
	"fmt"
	"log"
	"time"

	"gin-gonic/helper"
	"gin-gonic/modules"
	"gin-gonic/modules/books"
	"gin-gonic/modules/users"
	"gin-gonic/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// helper.()
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	if config.LOG_FILE == "on" {
		helper.SetupLogOutput()
	}

	app := gin.Default()
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.ALLOW_ORIGIN)
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API service is running",
		})
	})
	db := helper.OpenDb(config.DB, config.Schema, "v1")
	if db == nil {
		log.Fatal("Failed to Connect to database")

	}
	db = db.Debug()
	//connect NATS
	helper.ConnectNats(config.NatsUrl)
	if helper.NatsConn != nil{
		defer helper.NatsConn.Close()
	}
	// jalanin worker book
	bookService := books.NewBookService(db)
	books.StartWorker(bookService)

	seedAdmin()

	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS public").Error; err != nil {
		log.Fatal("Failed to create schema:", err)
	}

	if err := db.Exec("SET search_path TO public").Error; err != nil {
		log.Fatal("Failed to set search path:", err)
	}

	// // Middleware untuk logging dan recovery
	// app.Use(gin.Logger())
	// app.Use(gin.Recovery())

	versionRunner := modules.NewVersion(config, app, db, "api/v1")
	versionRunner.Run()

	fmt.Printf("Server starting on port %s\n", config.AppPort)
	app.Run(config.AppPort)
}

func seedAdmin() {
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	adminEmail := config.ADMIN_EMAIL
	adminPassword := config.AdminPassword
	adminName := config.AdminName

	// validasi env tidak diset
	if adminEmail == "" || adminPassword == "" {
		fmt.Println("seeding skipped: ADMIN_EMAIL or ADMIN_PASSWORD not found in .env")
	}
	var count int64

	//cek apakah admin sudah ada
	helper.DB.Model(&users.User{}).Where("role =?", "admin").Count(&count)

	if count == 0 {
		fmt.Println("No admin found. Creating admin from environment variable")
		//hash password dari env
		hashedPassword, err := utils.HashPassword(adminPassword)
		if err != nil {
			fmt.Printf("Error hashing password : %v\n", err)
			return
		}
		admin := users.User{
			Name:      adminName,
			Email:     adminEmail,
			Password:  hashedPassword,
			Address:   "System Administrator",
			Role:      "admin",
			BornDate:  time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := helper.DB.Create(&admin).Error; err != nil {
			fmt.Printf("failed to create admin :%v \n", err)
		} else {
			fmt.Println("admin account seeded succesfully")
		}

	} else {
		fmt.Println("Admin account check: OK (Admin already exists")
	}

}
