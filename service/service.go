package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/SeijiOmi/gin-tamplate/db"
	"github.com/SeijiOmi/gin-tamplate/entity"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Behavior ユーザサービスを提供するメソッド群
type Behavior struct{}

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
	db := db.GetDB()

	hash, err := createHashPassword(inputUser.Password)
	createUser.Password = hash
	if err != nil {
		return createUser, err
	}

	if err := db.Create(&createUser).Error; err != nil {
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

	user, err = b.GetByID(strconv.Itoa(id))
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
