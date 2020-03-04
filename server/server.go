package server

import (
	"github.com/SeijiOmi/gin-tamplate/controller"
	"github.com/gin-gonic/gin"
)

// Init is initialize server
func Init() {
	r := router()
	r.Run(":8080")
}

func router() *gin.Engine {
	r := gin.Default()

	u := r.Group("/users")
	{
		u.GET("", controller.Index)
		u.GET("/:id", controller.Show)
		u.POST("", controller.Create)
		u.PUT("/:id", controller.Update)
		u.DELETE("/:id", controller.Delete)

		u.POST("/login", controller.Login)
		u.POST("/auth", controller.Auth)
	}

	return r
}
