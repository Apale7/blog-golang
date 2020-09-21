package controller

import (
	"blog-gin/common"
	"blog-gin/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errors.WithStack(err))
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密错误"})
		return
	}
	log.Println(string(hashedPassword))
	newUser := model.User{
		Username: username,
		Password: string(hashedPassword),
	}
	db.Create(&newUser)
	ctx.JSON(200, gin.H{"msg": "注册成功"})
}

func Login(ctx *gin.Context)  {
	db := common.GetDb()
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
	var user model.User
	db.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户名或密码错误"})
		return
	}
	log.Println(user.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password));err!= nil {
		log.Error(errors.WithStack(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户名或密码错误"})
		return
	}
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Error(errors.WithStack(err))
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg": "登录成功",
	})
}

func UserInfo(ctx *gin.Context)  {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code":200,"data": gin.H{"user": user}})
}

func isUsernameExist(db *gorm.DB, username string) bool {
	var user model.User
	db.Where("username = ?", username).First(&user)
	return user.ID != 0
}
