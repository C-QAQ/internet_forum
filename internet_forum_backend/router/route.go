package router

import (
	"github.com/gin-gonic/gin"
	"internet_forum/controller"
	"internet_forum/logger"
	"internet_forum/middlewares"
	"net/http"
)

// Setup 路由
func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	//注册路由
	{
		v1.POST("/signup", controller.SignUpHandler)
		v1.POST("/login", controller.LoginHandler)
	}
	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/posts/", controller.GetPostListHandler)

		v1.POST("/vote", controller.PostVoteController)
	}

	v2 := r.Group("api/v2")
	//注册路由
	{
		v2.POST("/signup", controller.SignUpHandler)
		v2.POST("/login", controller.LoginHandler)
	}
	v2.Use(middlewares.JWTAuthMiddleware())
	{
		v2.GET("/community", controller.CommunityHandler)
		v2.GET("/community/:id", controller.CommunityDetailHandler)

		v2.POST("/post", controller.CreatePostHandler)
		v2.GET("/post/:id", controller.GetPostDetailHandler)
		v2.GET("/posts/", controller.GetPostListHandlerV2)

		v2.POST("/vote", controller.PostVoteController)
	}

	r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
