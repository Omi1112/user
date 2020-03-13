package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SeijiOmi/gin-tamplate/entity"
	"github.com/SeijiOmi/gin-tamplate/service"
)

// Index action: GET /users
func Index(c *gin.Context) {
	var b service.Behavior
	p, err := b.GetAll()

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
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
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(201, createdUser)
	}
}

// Show action: GET /users/:id
func Show(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior
	p, err := b.GetByID(id)

	if err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
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
		c.AbortWithStatus(400)
		fmt.Println(err)
	} else {
		c.JSON(200, p)
	}
}

// Delete action: DELETE /users/:id
func Delete(c *gin.Context) {
	id := c.Params.ByName("id")
	var b service.Behavior

	if err := b.DeleteByID(id); err != nil {
		c.AbortWithStatus(403)
		fmt.Println(err)
	} else {
		c.JSON(204, gin.H{"id #" + id: "deleted"})
	}
}

// Login action: POST /users/login
func Login(c *gin.Context) {
	var inputUser entity.User
	if err := bindJSON(c, &inputUser); err != nil {
		return
	}

	var b service.Behavior
	auth, err := b.LoginAuth(inputUser)
	if err != nil {
		c.AbortWithStatus(403)
		fmt.Println(err)
	} else {
		c.JSON(200, auth)
	}
}

// Auth action: POST /users/Auth
func Auth(c *gin.Context) {
	var b service.Behavior
	user, err := b.TokenAuth(c)
	if err != nil {
		c.AbortWithStatus(403)
		fmt.Println(err)
	} else {
		c.JSON(201, user)
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
