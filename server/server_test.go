package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SeijiOmi/gin-tamplate/db"
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
	router := router()

	req := httptest.NewRequest("GET", "/users", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	// assert.Equal(t, helloMessage, rec.Body.String())
}

// E2Eテストサンプル
func TestRouter(t *testing.T) {
	router := router()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL+"/users", nil)
	fmt.Println(testServer.URL + "/users")
	fmt.Println(req)

	resp, err := client.Do(req)
	fmt.Println(resp)
	fmt.Println(err)

	// respBody, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// assert.Equal(t, helloMessage, string(respBody))
}
