package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"

	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type env struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

func DB() *gorm.DB {
	godotenv.Load()
	loaded_env := env{
		DB_USERNAME: os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}

	// // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	dsn := "sqlserver://" + loaded_env.DB_USERNAME + ":" + loaded_env.DB_PASSWORD + "@" + loaded_env.DB_HOST + ":" + loaded_env.DB_PORT + "?" + "database=" + loaded_env.DB_NAME
	// fmt.Println(dsn)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func Migrate() {
	db := DB()
	// db.AutoMigrate(&UserDB{})
	// db.Table("MJM").AutoMigrate(&MJMDB{})
	db.Table("SFLA").AutoMigrate(&SFLADB{})
	db.Table("RM").AutoMigrate(&RMDB{})
}
