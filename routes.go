package main

import (
	"blog-gin/controller"
	"github.com/gin-gonic/gin"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.POST("/api/user/register", controller.Register)
	return r
}