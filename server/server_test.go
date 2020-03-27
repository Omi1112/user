package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/SeijiOmi/user/db"
	"github.com/SeijiOmi/user/entity"
	"github.com/stretchr/testify/assert"
)

var client = new(http.Client)
var testServer *httptest.Server
var userDefault = entity.User{Name: "test", Email: "test@mail.co.jp", Password: "password"}
var tmpBasePointURL string
var tmpBasePostURL string

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() {
	tmpBasePointURL = os.Getenv("POINT_URL")
	tmpBasePostURL = os.Getenv("POST_URL")
	setTestURL()
	db.Init()
	router := router()
	testServer = httptest.NewServer(router)
}

func teardown() {
	testServer.Close()
	db.Close()
	os.Setenv("POINT_URL", tmpBasePointURL)
	os.Setenv("POST_URL", tmpBasePostURL)
}

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

func TestCreateDemouser(t *testing.T) {
	initUserTable()
	resp, _ := http.Post(testServer.URL+"/demo", "application/json", nil)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func initUserTable() {
	db := db.GetDB()
	var u entity.User
	db.Delete(&u)
}

func setTestURL() {
	os.Setenv("POINT_URL", "http://user-mock-point:3000")
	os.Setenv("POST_URL", "http://user-mock-post:3000")
}
