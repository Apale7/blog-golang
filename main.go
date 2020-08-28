package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	_, err := os.Open("go.mod1")
	if err != nil {
		err = errors.WithStack(err)
		log.Warning(err)
	}
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Blog":   "www.flysnow.org",
			"wechat": "flysnow_org",
		})
	})
	r.Run(":8080")

}
