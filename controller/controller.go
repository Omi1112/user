package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SeijiOmi/user/entity"
	"github.com/SeijiOmi/user/service"
)

// Index action: GET /users
func Index(c *gin.Context) {
	var b service.Behavior
	p, err := b.GetAll()

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// Create action: POST /users
func Create(c *gin.Context) {
	var inputUser entity.User
	if err := bindJSON(c, &inputUser); err != nil {
		return
	}
	if userEmailExist(inputUser.Email) {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println("Email Exist")
		return
	}

	var b service.Behavior
	createdUser, err := b.CreateModel(inputUser)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, createdUser)
	}
}

// Show action: GET /users/:id
func Show(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	p, err := b.GetByID(id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// Update action: PUT /users/:id
func Update(c *gin.Context) {
	id := c.Params.ByName("id")
	var inputUser entity.User
	if err := bindJSON(c, &inputUser); err != nil {
		return
	}

	var b service.Behavior
	p, err := b.UpdateByID(id, inputUser)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, p)
	}
}

// Delete action: DELETE /users/:id
func Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior

	if err := b.DeleteByID(id); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, gin.H{"id #" + id: "deleted"})
	}
}

// Login action: POST /Auth
func Login(c *gin.Context) {
	inputUser := struct {
		Email    string
		Password string
	}{}
	if err := bindJSON(c, &inputUser); err != nil {
		return
	}

	var b service.Behavior
	auth, err := b.LoginAuth(inputUser.Email, inputUser.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, auth)
	}
}

// Auth action: GET /auth/:id
func Auth(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	user, err := b.TokenAuth(id)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

// CreateDemo action: POST /demo/:id
func CreateDemo(c *gin.Context) {
	var b service.Behavior
	demoUserAuth, err := b.CreateDemoData()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusCreated, demoUserAuth)
	}
}

func bindJSON(c *gin.Context, data interface{}) error {
	if err := c.BindJSON(data); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		fmt.Println("bind JSON err")
		fmt.Println(err)
		return err
	}
	return nil
}

func userEmailExist(email string) bool {
	var b service.Behavior
	user, _ := b.GetUserByEmail(email)
	if user.ID == 0 {
		return false
	}
	return true
}
