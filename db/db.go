package db

import (
	"fmt"
	"os"

	"github.com/SeijiOmi/gin-tamplate/entity"
	"github.com/jinzhu/gorm"

	// mysql呼び出し用の設定
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
	err error
)

// Init is initialize db from main function
func Init() {
	DBMS := "mysql"
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	PROTOCOL := "tcp(" + os.Getenv("DB_ADDRESS") + ")"
	DBNAME := os.Getenv("DB_NAME")
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true"
	fmt.Println(CONNECT)
	_db, err := gorm.Open(DBMS, CONNECT)
	db = _db
	if err != nil {
		panic(err)
	}

	autoMigration()
}

// GetDB is called in models
func GetDB() *gorm.DB {
	return db
}

// Close is closing db
func Close() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func autoMigration() {
	db.AutoMigrate(&entity.User{})
}
