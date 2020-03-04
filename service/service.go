package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/SeijiOmi/gin-tamplate/db"
	"github.com/SeijiOmi/gin-tamplate/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Behavior procides user's behavior
type Behavior struct{}

const (
	secret = "2FMd5FNSqS/nW2wWJy5S3ppjSHhUnLt8HuwBkTD6HqfPfBBDlykwLA=="
)

// GetAll is get all User
func (b Behavior) GetAll() ([]entity.User, error) {
	db := db.GetDB()
	var u []entity.User

	if err := db.Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

// CreateModel is create User model
func (b Behavior) CreateModel(inputUser entity.User) (entity.User, error) {
	var createUser entity.User
	db := db.GetDB()

	// パスワードのハッシュを生成
	hash, err := bcrypt.GenerateFromPassword([]byte(inputUser.Password), 10)
	inputUser.Password = string(hash)
	if err != nil {
		log.Fatal(err)
		return createUser, err
	}

	if err := db.Create(&createUser).Error; err != nil {
		return createUser, err
	}

	return createUser, nil
}

// GetByID is get a User
func (b Behavior) GetByID(id string) (entity.User, error) {
	db := db.GetDB()
	var u entity.User

	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// UpdateByID is update a User
func (b Behavior) UpdateByID(id string, c *gin.Context) (entity.User, error) {
	db := db.GetDB()
	var u entity.User

	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}

	if err := c.BindJSON(&u); err != nil {
		return u, err
	}

	db.Save(&u)

	return u, nil
}

// DeleteByID is delete a User
func (b Behavior) DeleteByID(id string) error {
	db := db.GetDB()
	var u entity.User

	if err := db.Where("id = ?", id).Delete(&u).Error; err != nil {
		return err
	}

	return nil
}

// LoginAuth ログイン認証を行い認証トークンを発行する
func (b Behavior) LoginAuth(inputUser entity.User) (string, error) {
	// ユーザの取得
	var dbUser entity.User
	db := db.GetDB()
	if err := db.Where("email = ?", inputUser.Email).First(&dbUser).Error; err != nil {
		return "", err
	}

	// パスワードの確認
	err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password))
	if err != nil {
		return "", err
	}

	// トークンの作成
	token, err := createToken(dbUser)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Println(token)

	return token, err
}

// TokenAuth 認証トークンで承認を行い、ユーザ情報を返却するサービス
func (b Behavior) TokenAuth(c *gin.Context) (entity.User, error) {
	var auth entity.Auth
	var user entity.User
	if err := c.BindJSON(&auth); err != nil {
		return user, err
	}
	id, err := perthToken(auth.Token)
	if err != nil {
		return user, err
	}
	fmt.Println(string(id))

	user, err = b.GetByID(strconv.Itoa(id))
	if err != nil {
		return user, err
	}

	return user, nil
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

// perthToken は jwt トークンからidを取得する。
func perthToken(signedString string) (int, error) {
	var id int
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
	id = int(floatID)
	return id, nil
}
