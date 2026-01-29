package main

import (
	"log"
	"nats-subscriber/helper"
	"nats-subscriber/modules"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

func main() {
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	if config.LOG_FILE == "on" {
		helper.SetupLogOutput()
	}

	gin.SetMode(config.GIN_MODE)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.ALLOW_ORIGIN)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API Service Subscriber is running",
		})
	})

	// Connect DB
	db := helper.OpenDB(config.DB, config.SCHEMA, "v1")
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	// db = db.Debug() // Optional

	// Connect NATS
	servers := strings.Split(config.NatsServers, ",")
	natsUrl := strings.Join(servers, ",")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}

	nc, err := nats.Connect(natsUrl,
		nats.Timeout(100*time.Second),
		nats.ReconnectWait(5*time.Second),
		nats.MaxReconnects(100),
		nats.DisconnectHandler(func(nc *nats.Conn) {
			log.Println("Disconnected from NATS")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Println("Reconnected to NATS")
		}),
	)

	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Subscriber started, listening for messages...")

	setupDatabaseSchema(db, config.SCHEMA)

	m_nats := modules.NewModulesNats(config)
	m_nats.Run(nc, db)

	if err := r.Run(":" + config.PORT); err != nil {
		log.Fatal(err)
	}
}
func setupDatabaseSchema(db *gorm.DB, schema string) {
	if schema == "" {
		return
	}
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + schema).Error; err != nil {
		log.Fatal("Failed to create schema:", err)
	}
	if err := db.Exec("SET search_path TO " + schema).Error; err != nil {
		log.Fatal("Failed to set search path:", err)
	}
}
