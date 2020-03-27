package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"io/ioutil"
	"strings"
	"encoding/json"

	"github.com/SeijiOmi/user/db"
	"github.com/SeijiOmi/user/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmcvetta/napping"
	"golang.org/x/crypto/bcrypt"
)

// Behavior ユーザサービスを提供するメソッド群
type Behavior struct{}

var demoOtherUser = entity.User{Name: "", Email: "test@co.jp", Password: "password"}
var demoLoginUser = entity.User{Name: "", Email: "test@co.jp", Password: "password"}

const (
	secret = "2FMd5FNSqS/nW2wWJy5S3ppjSHhUnLt8HuwBkTD6HqfPfBBDlykwLA=="
)

// GetAll ユーザ全件を取得
func (b Behavior) GetAll() ([]entity.User, error) {
	db := db.GetDB()
	var u []entity.User

	if err := db.Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

// CreateModel ユーザ情報の生成
func (b Behavior) CreateModel(inputUser entity.User) (entity.User, error) {
	createUser := inputUser

	hash, err := createHashPassword(inputUser.Password)
	createUser.Password = hash
	if err != nil {
		return createUser, err
	}

	db := db.GetDB()

	if err := db.Create(&createUser).Error; err != nil {
		return createUser, err
	}

	token, err := createToken(createUser)
	if err != nil {
		b.DeleteByID(strconv.Itoa(int(createUser.ID)))
		return createUser, err
	}
	initPoint := 10000
	err = createPoint(initPoint, "登録ありがとうございます！初期ポイントを付与します。" , token)
	if err != nil {
		b.DeleteByID(strconv.Itoa(int(createUser.ID)))
		return createUser, err
	}

	return createUser, nil
}

// GetByID IDを元にユーザ1件を取得
func (b Behavior) GetByID(id string) (entity.User, error) {
	db := db.GetDB()
	var u entity.User

	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// GetUserByEmail emailを基にユーザを取得する。
func (b Behavior) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	db := db.GetDB()
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

// UpdateByID 指定されたidをinputUser通りに更新
func (b Behavior) UpdateByID(id string, inputUser entity.User) (entity.User, error) {
	db := db.GetDB()
	var findUser entity.User
	if err := db.Where("id = ?", id).First(&findUser).Error; err != nil {
		return findUser, err
	}

	updateUser := inputUser
	hash, err := createHashPassword(inputUser.Password)
	updateUser.Password = hash
	if err != nil {
		return updateUser, err
	}
	updateUser.ID = findUser.ID
	db.Save(&updateUser)

	return updateUser, nil
}

// DeleteByID 指定されたidを削除
func (b Behavior) DeleteByID(id string) error {
	db := db.GetDB()
	var u entity.User

	if err := db.Where("id = ?", id).Delete(&u).Error; err != nil {
		return err
	}

	return nil
}

// LoginAuth ログイン認証を行い認証トークンを発行
func (b Behavior) LoginAuth(email string, password string) (entity.Auth, error) {
	dbUser, err := b.GetUserByEmail(email)
	if err != nil {
		return entity.Auth{}, err
	}

	// パスワードの確認
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if err != nil {
		return entity.Auth{}, err
	}

	// トークンの作成
	token, err := createToken(dbUser)
	if err != nil {
		log.Fatal(err)
		return entity.Auth{}, err
	}
	fmt.Println(token)
	returnAuth := entity.Auth{Token: token, ID: dbUser.ID}

	return returnAuth, err
}

// TokenAuth 認証トークンで承認を行い、ユーザ情報を返却するサービス
func (b Behavior) TokenAuth(token string) (entity.User, error) {
	var user entity.User
	id, err := perthToken(token)
	if err != nil {
		return user, err
	}
	fmt.Println(string(id))

	user, err = b.GetByID(strconv.Itoa(int(id)))
	if err != nil {
		return user, err
	}

	return user, nil
}

func createHashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	hashPassword := string(hash)
	if err != nil {
		log.Fatal(err)
		return hashPassword, err
	}
	return hashPassword, nil
}

// Token 作成関数
func createToken(u entity.User) (string, error) {
	var err error

	// Token を作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": u.ID,
	})
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}

// perthToken jwt トークンからidを取得する。
func perthToken(signedString string) (uint, error) {
	var id uint
	token, err := jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return id, fmt.Errorf("%s is expired", signedString)
			}
			return id, fmt.Errorf("%s is invalid", signedString)
		}
		return id, fmt.Errorf("%s is invalid", signedString)
	}

	if token == nil {
		return 0, fmt.Errorf("not found token in %s", signedString)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("not found claims in %s", signedString)
	}

	floatID, ok := claims["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("not found id in %s", signedString)
	}
	id = uint(floatID)
	return id, nil
}

func createPoint(point int, comment string, token string) error {
	input := struct {
		Number int    `json:"number"`
		Comment  string `json:"comment"`
		Token  string `json:"token"`
	}{
		point,
		comment,
		token,
	}
	error := struct {
		Error string
	}{}

	baseURL := os.Getenv("POINT_URL")
	resp, err := napping.Post(baseURL+"/points", &input, nil, &error)

	if err != nil {
		return err
	}

	if resp.Status() == http.StatusBadRequest {
		return errors.New("token invalid")
	}

	return nil
}

// CreateDemoData デモデータを作成して、デモユーザーのログイン情報を返却する
func (b Behavior) CreateDemoData() (entity.Auth, error) {
	otherUserAuth, err := b.createUniqueDemoUser(demoOtherUser)
	if err != nil {
		return entity.Auth{}, err
	}

	loginUserAuth, err := b.createUniqueDemoUser(demoLoginUser)
	if err != nil {
		return entity.Auth{}, err
	}

	err = createOtherServiceData(otherUserAuth, loginUserAuth)
	if err != nil {
		return entity.Auth{}, err
	}

	return loginUserAuth, nil
}

func (b Behavior) createUniqueDemoUser(user entity.User) (entity.Auth, error) {
	maxID, err := getMaxUserID()
	if err != nil {
		return entity.Auth{}, err
	}
	uniqueDmoUser := user
	maxID++
	uniqueDmoUser.Name, err = getDemoUserName()
	if err != nil {
		return entity.Auth{}, err
	}
	uniqueDmoUser.Email = strconv.Itoa(maxID) + uniqueDmoUser.Email
	_, err = b.CreateModel(uniqueDmoUser)
	if err != nil {
		return entity.Auth{}, err
	}
	uniqueDmoUserAuth, err := b.LoginAuth(uniqueDmoUser.Email, uniqueDmoUser.Password)
	if err != nil {
		return entity.Auth{}, err
	}

	return uniqueDmoUserAuth, nil
}

func getMaxUserID() (int, error) {
	db := db.GetDB()
	rows, err := db.Table("users").Select("max(id) as maxID").Rows()
	defer rows.Close()
	if err != nil {
		return 0, err
	}
	var maxID int
	for rows.Next() {
		rows.Scan(&maxID)
	}

	return maxID, nil
}

func createOtherServiceData(otherUserAuth entity.Auth, loginUserAuth entity.Auth) error {
	demoData := []struct {
		body        string
		point       int
		token       string
		helperToken string
		doneToken   string
	}{
		{
			"お風呂入っている間子供を見ててくれませんか？",
			100,
			otherUserAuth.Token,
			loginUserAuth.Token,
			"",
		},
		{
			"模様替えの家具移動手伝って下さい！",
			500,
			otherUserAuth.Token,
			loginUserAuth.Token,
			otherUserAuth.Token,
		},
		{
			"コストコ会員の人、一緒に連れてってくれませんか？",
			100,
			otherUserAuth.Token,
			"",
			"",
		},
		{
			"背の高い人電球交換助けてくれませんか？",
			100,
			loginUserAuth.Token,
			otherUserAuth.Token,
			loginUserAuth.Token,
		},
		{
			"テレビがつきません！！詳しい人いませんか？",
			200,
			loginUserAuth.Token,
			otherUserAuth.Token,
			"",
		},
		{
			"ベットの組み立て手伝ってください！！",
			100,
			loginUserAuth.Token,
			"",
			"",
		},
	}

	for _, data := range demoData {
		postID, err := createPost(
			data.body,
			data.point,
			data.token)
		if err != nil {
			return err
		}

		if data.helperToken != "" {
			err := setHelperPost(postID, data.helperToken)
			if err != nil {
				return err
			}
		}

		if data.doneToken != "" {
			err := donePost(postID, data.doneToken)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createPost(body string, point int, token string) (int, error) {
	input := struct {
		Body  string `json:"body"`
		Point int    `json:"point"`
		Token string `json:"token"`
	}{
		body,
		point,
		token,
	}
	response := struct {
		ID int
	}{}
	error := struct {
		Error string
	}{}

	baseURL := os.Getenv("POST_URL")
	resp, err := napping.Post(baseURL+"/posts", &input, &response, &error)

	if err != nil {
		return 0, err
	}

	if resp.Status() == http.StatusBadRequest {
		return 0, errors.New("token invalid")
	}

	return response.ID, nil
}

func setHelperPost(id int, token string) error {
	input := struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	}{
		id,
		token,
	}
	error := struct {
		Error string
	}{}

	baseURL := os.Getenv("POST_URL")
	resp, err := napping.Post(baseURL+"/helper", &input, nil, &error)

	if err != nil {
		return err
	}

	if resp.Status() == http.StatusBadRequest {
		return errors.New("token invalid")
	}

	return nil
}

func donePost(id int, token string) error {
	input := struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	}{
		id,
		token,
	}
	error := struct {
		Error string
	}{}

	baseURL := os.Getenv("POST_URL")
	resp, err := napping.Post(baseURL+"/done", &input, nil, &error)

	if err != nil {
		return err
	}

	if resp.Status() == http.StatusBadRequest {
		return errors.New("token invalid")
	}

	return nil
}

func getDemoUserName() (string, error) {
	res, err := http.Get("https://green.adam.ne.jp/roomazi/cgi-bin/randomname.cgi?n=1")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	jsonpStr := string(byteArray)
	jsonStr := jsonpStr[strings.Index(jsonpStr, "(")+1 : strings.Index(jsonpStr, ")")]

	response := struct {
		Err  int
		Name [][]string
	}{}
	json.Unmarshal([]byte(jsonStr), &response)

	return response.Name[0][0], nil
}
