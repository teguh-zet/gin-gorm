package modules

import (
	"log"
	"nats-subscriber/helper"
	"nats-subscriber/modules/loans"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type ModulesNats interface {
	Run(nc *nats.Conn, db *gorm.DB)
}

type modulesNats struct {
	config helper.Config
}

func NewModulesNats(config helper.Config) ModulesNats {
	return &modulesNats{config: config}
}

func (m *modulesNats) Run(nc *nats.Conn, db *gorm.DB) {
	log.Println("Modules Nats Started")

	svr := make(chan string)

	// Check AutoMigrate
	autoMigrate := m.config.AUTO_MIGRATE == "Y" || m.config.AUTO_MIGRATE == "on" || m.config.AUTO_MIGRATE == "true"

	// Init Loan Server
	loanServer := loans.NewLoanNatsServer(db, nc, autoMigrate)
	
	// Init blocks, so run in goroutine
	go loanServer.Init(svr)
	
	// We could handle svr messages here if we wanted to track status
	go func() {
		for s := range svr {
			log.Printf("Module %s stopped/initialized", s)
		}
	}()
}
