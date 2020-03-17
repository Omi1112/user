package service

import (
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/SeijiOmi/gin-tamplate/db"
	"github.com/SeijiOmi/gin-tamplate/entity"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

/*
	テストの前準備
*/

var client = new(http.Client)
var userDefault = entity.User{Name: "test", Email: "test@co.jp", Password: "password"}
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
	initUserTable()
}

// テスト実施後共通処理
func teardown() {
	db.Close()
	os.Setenv("POINT_URL", tmpBasePointURL)
}

/*
	ここからが個別のテスト実装
*/

func TestGetAll(t *testing.T) {
	initUserTable()
	createDefaultUser()
	createDefaultUser()

	var b Behavior
	users, err := b.GetAll()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(users), 2)
}

func TestCreateModel(t *testing.T) {
	var b Behavior
	user, err := b.CreateModel(userDefault)

	assert.Equal(t, nil, err)
	assert.Equal(t, userDefault.Name, user.Name)
	assert.Equal(t, userDefault.Email, user.Email)
	assert.NotEqual(t, userDefault.Password, user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDefault.Password))
	assert.Equal(t, nil, err)
}

func TestGetByIDExists(t *testing.T) {
	user := createDefaultUser()
	var b Behavior
	user, err := b.GetByID(strconv.Itoa(int(user.ID)))

	assert.Equal(t, nil, err)
	assert.Equal(t, userDefault.Name, user.Name)
	assert.Equal(t, userDefault.Email, user.Email)
}

func TestGetByIDNotExists(t *testing.T) {
	var b Behavior
	user, err := b.GetByID(string(userDefault.ID))

	assert.NotEqual(t, nil, err)
	var nilUser entity.User
	assert.Equal(t, nilUser, user)
}

func TestUpdateByIDExists(t *testing.T) {
	user := createDefaultUser()

	updateUser := entity.User{Name: "not", Email: "not@co.jp", Password: "notpassword"}

	var b Behavior
	user, err := b.UpdateByID(strconv.Itoa(int(user.ID)), updateUser)

	assert.Equal(t, nil, err)
	assert.Equal(t, updateUser.Name, user.Name)
	assert.Equal(t, updateUser.Email, user.Email)
	assert.NotEqual(t, updateUser.Password, user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateUser.Password))
	assert.Equal(t, nil, err)
}

func TestUpdateByIDNotExists(t *testing.T) {
	user := createDefaultUser()

	updateUser := entity.User{Name: "not", Email: "not@co.jp", Password: "notpassword"}

	var b Behavior
	user, err := b.UpdateByID("0", updateUser)

	assert.NotEqual(t, nil, err)
	var nilUser entity.User
	assert.Equal(t, nilUser, user)
}

func TestDeleteByIDExists(t *testing.T) {
	user := createDefaultUser()

	db := db.GetDB()
	var beforeCount int
	db.Table("users").Count(&beforeCount)

	var b Behavior
	err := b.DeleteByID(strconv.Itoa(int(user.ID)))

	var afterCount int
	db.Table("users").Count(&afterCount)

	assert.Equal(t, nil, err)
	assert.Equal(t, beforeCount-1, afterCount)
}

func TestDeleteByIDNotExists(t *testing.T) {
	initUserTable()
	createDefaultUser()

	db := db.GetDB()
	var beforeCount int
	db.Table("users").Count(&beforeCount)

	var b Behavior
	err := b.DeleteByID("0")

	var afterCount int
	db.Table("users").Count(&afterCount)

	assert.Equal(t, nil, err)
	assert.Equal(t, beforeCount, afterCount)
}

func createDefaultUser() entity.User {
	db := db.GetDB()
	user := userDefault
	db.Create(&user)
	return user
}

func initUserTable() {
	db := db.GetDB()
	var u entity.User
	db.Delete(&u)
}
