package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DB_USERNAME = "root"
	DB_PASSWORD = "password"
	DB_HOST     = "127.0.0.1"
	DB_PORT     = "3306"
	DB_NAME     = "db"
)

func DB() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func Migrate() {
	db := DB()
	// db.AutoMigrate(&UserDB{})
	db.Table("MJM").AutoMigrate(&MJMDB{})
	db.Table("SFLA").AutoMigrate(&SFLADB{})
	db.Table("RM").AutoMigrate(&RMDB{})
}
