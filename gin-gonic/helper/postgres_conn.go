package helper

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func OpenDb(conn string, schm string, ver string) *gorm.DB {
	dsn := conn
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   schm + ".",
			SingularTable: true,
		},
	})

	if err != nil {
		panic("Failed to connect database")
	}
	DB = db
	return db

}

//manual
// var DB *gorm.DB

// func ConnectDatabase() {
// 	var errConnection error

// 	dbHost := GetConfig("DB_HOST")
// 	dbDriver := GetConfig("DB_DRIVER")
// 	dbName := GetConfig("DB_NAME")
// 	dbUser := GetConfig("DB_USER")
// 	dbPass := GetConfig("DB_PASSWORD")
// 	dbPort := GetConfig("DB_PORT")
// 	if dbDriver == "mysql" {
// 		dsnMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
// 		DB, errConnection = gorm.Open(mysql.Open(dsnMysql), &gorm.Config{})
// 	}

// 	if dbDriver == "postgre" {
// 		dsnp := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, dbUser, dbPass, dbName, dbPort)
// 		DB, errConnection = gorm.Open(postgres.Open(dsnp), &gorm.Config{})
// 	}

// 	if errConnection != nil {
// 		panic("can't connect to database")
// 	}
// }

// func CloseDB(db *gorm.DB){
// 	sqlDB, err := db.DB()
// 	if err!= nil{
// 		panic("Failed to close database")
// 	}
// 	sqlDB.Close()
// }
