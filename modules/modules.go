package modules

import (
	"gin-gonic/helper"
	"gin-gonic/modules/books"
	"gin-gonic/modules/loans"
	"gin-gonic/modules/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Versions interface {
	Run()
}

type versions struct {
	mainServer *gin.Engine
	db         *gorm.DB
	version    string
}

func NewVersion(config helper.Config, router *gin.Engine, db *gorm.DB, version string) Versions {
	return &versions{
		mainServer: router,
		db:         db,
		version:    version,
	}
}

func (s *versions) Run() {
	apiRoutes := s.mainServer.Group("/")

	userServer := users.NewUserServer(apiRoutes, s.db, s.version)
	userServer.Init()

	bookServer := books.NewBookServer(apiRoutes, s.db, s.version)
	bookServer.Init()

	loanServer := loans.NewLoanServer(apiRoutes, s.db, s.version)
	loanServer.Init()
}
