package main

import (
	"github.com/gin-gonic/gin"
	_ "blog-gin/db"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Blog":   "www.flysnow.org",
			"wechat": "flysnow_org",
		})
	})
	r.Run(":8080")

}
