package apm_rawat_inap

import (
	"log"

	"gin-gonic/helper"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ApmRawatInapServer struct {
	router  *gin.RouterGroup
	db      *gorm.DB
	nc      *nats.Conn
	version string
}

func NewApmRawatInapServer(router *gin.RouterGroup, db *gorm.DB, nc *nats.Conn, version string) *ApmRawatInapServer {
	return &ApmRawatInapServer{router: router, db: db, nc: nc, version: version}
}

func (s *ApmRawatInapServer) Init() {
	config, err := helper.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	if config.AUTO_MIGRATE == "Y" {
		if err := s.db.AutoMigrate(&ApmRawatInap{}); err != nil {
			log.Printf("Failed to auto migrate ApmRawatInap: %v", err)
		}
	}

	service := NewApmRawatInapService(s.db, s.nc)
	controller := NewApmRawatInapController(service)

	routes := s.router.Group("/" + s.version + "/apm-rawat-inap")
	routes.POST("", controller.Create)
	routes.POST("/:id/publish", controller.Publish)
}
