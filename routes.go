package main

import (
	"blog-gin/controller"
	"blog-gin/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.POST("/api/user/register", controller.Register)
	r.POST("/api/user/login", controller.Login)
	r.POST("/api/user/logout", controller.Logout)
	r.GET("/api/user/whoami",middleware.LoginCheck(), controller.WhoAmI)
	r.POST("/api/blog/create", middleware.LoginCheck(), controller.CreateBlog)
	r.POST("/api/blog/update", middleware.LoginCheck(), controller.UpdateBlog)
	r.POST("/api/blog/delete", middleware.LoginCheck(), controller.DeleteBlog)
	r.GET("/api/user/info", middleware.LoginCheck(), controller.UserInfo)
	r.GET("/api/blog/list",  controller.BlogList)
	r.GET("/api/blog/content",  controller.BlogContent)
	return r
}
