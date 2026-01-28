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
	"gin-gonic/websocket"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 1. Load Configuration & Logger
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	if config.LOG_FILE == "on" {
		helper.SetupLogOutput()
	}

	// 2. Initialize Database
	db := helper.OpenDb(config.DB, config.Schema, "v1")
	if db == nil {
		log.Fatal("Failed to connect to database")
	}
	// db = db.Debug() // Uncomment jika ingin mode debug query

	// 3. Connect NATS (PENTING: Harus sebelum WebSocket Manager)
	helper.ConnectNats(config.NatsUrl)
	if helper.NatsConn != nil {
		defer helper.NatsConn.Close()
	}

	// 4. Initialize WebSocket Manager
	// Manager ini butuh NATS yang sudah terkoneksi untuk subscribe topic
	wsManager := websocket.NewManager()
	go wsManager.Run()

	// 5. Setup Gin Engine & Middleware
	app := gin.Default()
	app.Use(CORSMiddleware(config.ALLOW_ORIGIN))

	// 6. Define Routes
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API service is running"})
	})

	// Route khusus WebSocket
	app.GET("/ws", func(c *gin.Context) {
		websocket.ServeWS(wsManager, c)
	})

	// Setup API Versioning & Modules
	versionRunner := modules.NewVersion(config, app, db, "api/v1")
	versionRunner.Run()

	// 7. Background Workers & Database Setup
	bookService := books.NewBookService(db)
	books.StartWorker(bookService)

	setupDatabaseSchema(db)
	seedAdmin(db, config) // Pass db & config agar efisien

	// 8. Start Server
	fmt.Printf("Server starting on port %s\n", config.AppPort)
	if err := app.Run(config.AppPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// CORSMiddleware memisahkan logika CORS agar main function lebih bersih
func CORSMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// setupDatabaseSchema menangani konfigurasi schema Postgres
func setupDatabaseSchema(db *gorm.DB) {
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS public").Error; err != nil {
		log.Fatal("Failed to create schema:", err)
	}
	if err := db.Exec("SET search_path TO public").Error; err != nil {
		log.Fatal("Failed to set search path:", err)
	}
}

// seedAdmin sekarang menerima DB dan Config sebagai parameter
func seedAdmin(db *gorm.DB, config helper.Config) {
	adminEmail := config.ADMIN_EMAIL
	adminPassword := config.AdminPassword
	adminName := config.AdminName

	if adminEmail == "" || adminPassword == "" {
		fmt.Println("Seeding skipped: ADMIN_EMAIL or ADMIN_PASSWORD not found in .env")
		return
	}

	var count int64
	if err := db.Model(&users.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		fmt.Printf("Error checking admin existence: %v\n", err)
		return
	}

	if count == 0 {
		fmt.Println("No admin found. Creating admin from environment variable...")
		
		hashedPassword, err := utils.HashPassword(adminPassword)
		if err != nil {
			fmt.Printf("Error hashing password: %v\n", err)
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

		if err := db.Create(&admin).Error; err != nil {
			fmt.Printf("Failed to create admin: %v\n", err)
		} else {
			fmt.Println("Admin account seeded successfully")
		}
	} else {
		fmt.Println("Admin account check: OK (Admin already exists)")
	}
}
// package main

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"gin-gonic/helper"
// 	"gin-gonic/modules"
// 	"gin-gonic/modules/books"
// 	"gin-gonic/modules/users"
// 	"gin-gonic/utils"
// 	"gin-gonic/websocket"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	// helper.()
// 	config, err := helper.LoadConfig(".")
// 	if err != nil {
// 		log.Fatal("cannot load config")
// 	}
// 	if config.LOG_FILE == "on" {
// 		helper.SetupLogOutput()
// 	}
// 	//WEBSOCKET
// 	wsManager := websocket.NewManager() 
//     go wsManager.Run()


// 	app := gin.Default()
// 	app.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", config.ALLOW_ORIGIN)
// 		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	})
// 	app.GET("/ws", func(c *gin.Context) {
//         // 3. KITA OPER MANAGER INI KE FUNCTION SERVEWS
//         websocket.ServeWS(wsManager, c)
//     })
// 	app.GET("/", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "API service is running",
// 		})
// 	})
// 	db := helper.OpenDb(config.DB, config.Schema, "v1")
// 	if db == nil {
// 		log.Fatal("Failed to Connect to database")

// 	}
// 	db = db.Debug()
// 	//connect NATS
// 	helper.ConnectNats(config.NatsUrl)
// 	if helper.NatsConn != nil{
// 		defer helper.NatsConn.Close()
// 	}
// 	// jalanin worker book
// 	bookService := books.NewBookService(db)
// 	books.StartWorker(bookService)

// 	seedAdmin()

// 	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS public").Error; err != nil {
// 		log.Fatal("Failed to create schema:", err)
// 	}

// 	if err := db.Exec("SET search_path TO public").Error; err != nil {
// 		log.Fatal("Failed to set search path:", err)
// 	}

// 	// // Middleware untuk logging dan recovery
// 	// app.Use(gin.Logger())
// 	// app.Use(gin.Recovery())

// 	versionRunner := modules.NewVersion(config, app, db, "api/v1")
// 	versionRunner.Run()

// 	fmt.Printf("Server starting on port %s\n", config.AppPort)
// 	app.Run(config.AppPort)
// }

// func seedAdmin() {
// 	config, err := helper.LoadConfig(".")
// 	if err != nil {
// 		log.Fatal("cannot load config")
// 	}

// 	adminEmail := config.ADMIN_EMAIL
// 	adminPassword := config.AdminPassword
// 	adminName := config.AdminName

// 	// validasi env tidak diset
// 	if adminEmail == "" || adminPassword == "" {
// 		fmt.Println("seeding skipped: ADMIN_EMAIL or ADMIN_PASSWORD not found in .env")
// 	}
// 	var count int64

// 	//cek apakah admin sudah ada
// 	helper.DB.Model(&users.User{}).Where("role =?", "admin").Count(&count)

// 	if count == 0 {
// 		fmt.Println("No admin found. Creating admin from environment variable")
// 		//hash password dari env
// 		hashedPassword, err := utils.HashPassword(adminPassword)
// 		if err != nil {
// 			fmt.Printf("Error hashing password : %v\n", err)
// 			return
// 		}
// 		admin := users.User{
// 			Name:      adminName,
// 			Email:     adminEmail,
// 			Password:  hashedPassword,
// 			Address:   "System Administrator",
// 			Role:      "admin",
// 			BornDate:  time.Now(),
// 			CreatedAt: time.Now(),
// 			UpdatedAt: time.Now(),
// 		}
// 		if err := helper.DB.Create(&admin).Error; err != nil {
// 			fmt.Printf("failed to create admin :%v \n", err)
// 		} else {
// 			fmt.Println("admin account seeded succesfully")
// 		}

// 	} else {
// 		fmt.Println("Admin account check: OK (Admin already exists")
// 	}

// }
