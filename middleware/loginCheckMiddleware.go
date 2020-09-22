package middleware

import (
	"blog-gin/common"
	"blog-gin/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func LoginCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "权限不足",
			})
			ctx.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString, false)
		if err != nil || !token.Valid{
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "令牌已过期，请刷新"})
			ctx.Abort()
			return
		}
		userId:=claims.UserId
		db := common.GetDb()
		var user model.User
		db.First(&user, userId)
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
	}
}
