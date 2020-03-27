package db

import (
	"fmt"
	"os"

	"github.com/SeijiOmi/user/entity"
	"github.com/jinzhu/gorm"

	// mysql呼び出し用の設定
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db  *gorm.DB
	err error
)

// Init DB初期設定
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

// GetDB DBアクセサ取得
func GetDB() *gorm.DB {
	return db
}

// Close DB接続終了
func Close() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func autoMigration() {
	db.AutoMigrate(&entity.User{})
}
