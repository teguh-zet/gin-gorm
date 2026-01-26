package helpers

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var errConnection error

	dbHost := GetConfig("DB_HOST")
	dbDriver := GetConfig("DB_DRIVER")
	dbName := GetConfig("DB_NAME")
	dbUser := GetConfig("DB_USER")
	dbPass := GetConfig("DB_PASSWORD")
	dbPort := GetConfig("DB_PORT")
	if dbDriver == "mysql" {
		dsnMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
		DB, errConnection = gorm.Open(mysql.Open(dsnMysql), &gorm.Config{})
	}

	if dbDriver == "postgre" {
		dsnp := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPass, dbName, dbPort)
		DB, errConnection = gorm.Open(postgres.Open(dsnp), &gorm.Config{})
	}

	if errConnection != nil {
		panic("can't connect to database")
	}
}
