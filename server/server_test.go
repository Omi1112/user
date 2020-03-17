package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
var testServer *httptest.Server
var userDefault = entity.User{Name: "test", Email: "test@mail.co.jp", Password: "password"}
var tmpBasePointURL string

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
	tmpBasePointURL = os.Getenv("POINT_URL")
	os.Setenv("POINT_URL", "http://user-mock-point:3000")
	db.Init()
	router := router()
	testServer = httptest.NewServer(router)
}

// テスト実施後共通処理
func teardown() {
	testServer.Close()
	db.Close()
	os.Setenv("POINT_URL", tmpBasePointURL)
}

/*
	ここからが個別のテスト実装
*/

func TestUserCreateSuccessValid(t *testing.T) {
	initUserTable()
	input, _ := json.Marshal(userDefault)
	resp, _ := http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUserCreateRequireErrValid(t *testing.T) {
	initUserTable()
	inputUser := userDefault
	inputUser.Name = ""
	input, _ := json.Marshal(inputUser)
	resp, _ := http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	initUserTable()
	inputUser = userDefault
	inputUser.Email = ""
	input, _ = json.Marshal(inputUser)
	resp, _ = http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	initUserTable()
	inputUser = userDefault
	inputUser.Password = ""
	input, _ = json.Marshal(inputUser)
	resp, _ = http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUserCreateFormatErrValid(t *testing.T) {
	initUserTable()
	inputUser := userDefault
	inputUser.Email = "test"
	input, _ := json.Marshal(inputUser)
	resp, _ := http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
func TestUserCreateAlphaNumErrValid(t *testing.T) {
	initUserTable()
	inputUser := userDefault
	inputUser.Password = "あいうえおかきくけこ"
	input, _ := json.Marshal(inputUser)
	resp, _ := http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUserCreateEmailUniqueErrValid(t *testing.T) {
	initUserTable()
	inputUser := userDefault
	input, _ := json.Marshal(inputUser)
	resp, _ := http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, _ = http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLoginSuccess(t *testing.T) {
	initUserTable()
	input, _ := json.Marshal(userDefault)
	resp, _ := http.Post(testServer.URL+"/users", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	loginUser := struct {
		Email    string
		Password string
	}{
		userDefault.Email,
		userDefault.Password,
	}

	input, _ = json.Marshal(loginUser)
	resp, _ = http.Post(testServer.URL+"/auth", "application/json", bytes.NewBuffer(input))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func initUserTable() {
	db := db.GetDB()
	var u entity.User
	db.Delete(&u)
}
