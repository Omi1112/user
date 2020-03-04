package service

import (
	"net/http"
	"os"
	"testing"

	"github.com/SeijiOmi/gin-tamplate/db"
	"github.com/SeijiOmi/gin-tamplate/entity"
	"github.com/stretchr/testify/assert"
)

/*
	テストの前準備
*/

var client = new(http.Client)

// テストを統括するテスト時には、これが実行されるイメージでいる。
func TestMain(m *testing.M) {
	// テスト実施前の共通処理（自作関数）
	setup()
	ret := m.Run()
	// テスト実施後の共通処理（自作関数）
	teardown()
	os.Exit(ret)
}

// テスト実施前共通処理
func setup() {
	db.Init()
	db := db.GetDB()
	// DB初期化
	var u entity.User
	db.Delete(&u)
}

// テスト実施後共通処理
func teardown() {
	db.Close()
}

/*
	ここからが個別のテスト実装
*/

// サーバー内部でのテストサンプル
func TestHelloHandler(t *testing.T) {
	assert.Equal(t, 1, 1)
}
