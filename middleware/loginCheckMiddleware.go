package middleware

import (
	"blog-gin/common"
	"blog-gin/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		db := common.GetDb()
		user := model.User{}
		err := db.Model(&model.User{}).Where("id=?",userID.(uint)).First(&user).Error
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "系统错误"})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
