package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func LoginCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid, err := c.Cookie("userid")
		if err != nil {
			log.Error(errors.WithStack(err))
		}
		if userid == "apale" {
			c.Next()
		}else {
			c.Abort()
			c.SetCookie("userid", "apale", 10, "/", "localhost", false, true)
			c.String(400, "登录失败")
		}
	}
}