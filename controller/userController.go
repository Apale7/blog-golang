package controller

import (
	"blog-gin/common"
	"blog-gin/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func Register(ctx *gin.Context) {
	db:=common.GetDb()
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if len(username) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户名不能为空"})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码长度不能少于6位"})
		return
	}
	if isUsernameExist(db, username) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户名已存在"})
		return
	}
	newUser := model.User{
		Username: username,
		Password: password,
	}
	common.Db.Create(&newUser)
	ctx.JSON(200, gin.H{"msg": "注册成功"})
}

func isUsernameExist(db *gorm.DB, username string) bool {
	var user model.User
	db.Where("username = ?", username).First(&user)
	return user.ID != 0
}
