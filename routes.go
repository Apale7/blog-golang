package main

import (
	"blog-gin/controller"
	"blog-gin/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.POST("/api/user/register", controller.Register)
	r.POST("/api/user/login", controller.Login)
	r.POST("/api/user/refresh", controller.Refresh)
	r.POST("/api/blog/create",middleware.LoginCheck(), controller.CreateBlog)
	r.POST("/api/blog/update",middleware.LoginCheck(), controller.UpdateBlog)
	r.GET("/api/user/info", middleware.LoginCheck(),controller.UserInfo)
	return r
}