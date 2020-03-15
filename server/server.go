package server

import (
	"github.com/SeijiOmi/gin-tamplate/controller"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Init is initialize server
func Init() {
	r := router()
	r.Run(":8080")
}

func router() *gin.Engine {
	r := gin.Default()

	// https://godoc.org/github.com/gin-gonic/gin#RouterGroup.Use
	r.Use(cors.New(cors.Config{
		// 許可したいHTTPメソッドの一覧
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
			"PUT",
			"DELETE",
		},
		// 許可したいHTTPリクエストヘッダの一覧
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"X-Requested-With",
			"Origin",
			"X-Csrftoken",
			"Content-Type",
			"Accept",
		},
		// 許可したいアクセス元の一覧
		AllowOrigins: []string{
			"*",
		},
	}))

	u := r.Group("/users")
	{
		u.GET("", controller.Index)
		u.GET("/:id", controller.Show)
		u.POST("", controller.Create)
		u.PUT("/:id", controller.Update)
		u.DELETE("/:id", controller.Delete)
	}

	a := r.Group("/auth")
	{
		a.GET("/:id", controller.Auth)
		a.POST("", controller.Login)
	}

	return r
}
