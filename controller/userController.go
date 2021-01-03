package controller

import (
	"blog-gin/common"
	"blog-gin/dto"
	"blog-gin/model"
	"blog-gin/response"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
)

func WhoAmI(c *gin.Context) {
	userInterface, exist := c.Get("user")
	if !exist {
		response.Fail(c, nil, "未登录")
		return
	}
	user := userInterface.(model.User)
	response.Success(c, gin.H{
		"nickname": user.Nickname,
		"username": user.Username,
		"id":       user.ID,
	}, "你是"+user.Nickname)
}

func Register(ctx *gin.Context) {
	db := common.GetDb()
	userInfo := struct {
		Username string
		Password string
		Nickname string
	}{}
	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	err = json.Unmarshal(bytes, &userInfo)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求参数错误")
		return
	}
	username, password, nickname := userInfo.Username, userInfo.Password, userInfo.Nickname
	if len(username) == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户名不能为空")
		return
	}
	if len([]byte(username)) > 32 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户名过长")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码长度不能少于6位")

		return
	}
	if isUsernameExist(db, username) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户名已存在")

		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errors.WithStack(err))
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	if len(nickname) == 0 {
		nickname = username
	}
	newUser := model.User{
		Username: username,
		Nickname: html.EscapeString(nickname),
		Password: string(hashedPassword),
	}
	db.Create(&newUser)
	response.Success(ctx, nil, "注册成功")
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	db := common.GetDb()
	userInfo := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := c.BindJSON(&userInfo); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数错误")
		return
	}
	username, password := userInfo.Username, userInfo.Password
	if len(username) == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户名或密码错误")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户名或密码错误")
		return
	}
	var user model.User
	db.Where("username = ?", username).Take(&user)
	if user.ID == 0 {
		response.Response(c, http.StatusBadRequest, 400, nil, "用户名或密码错误")
		return
	}
	//log.Println(user.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Error(errors.WithStack(err))
		response.Response(c, http.StatusBadRequest, 400, nil, "用户名或密码错误")
		return
	}

	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	response.Success(c, gin.H{
		"nickname": user.Nickname,
		"username": user.Username,
		"id":       user.ID,
	}, "登录成功")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user_id")
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete("user_id")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func Refresh(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		response.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		ctx.Abort()
		return
	}
	tokenString = tokenString[7:]
	token, claims, err := common.ParseToken(tokenString, true)
	if err != nil || !token.Valid {
		response.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		ctx.Abort()
		return
	}
	db := common.GetDb()
	var user model.User
	db.First(&user, claims.UserId)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnauthorized, 401, nil, "权限不足")
		ctx.Abort()
		return
	}
	accessToken, err := common.ReleaseToken(user, false)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	refreshToken, err := common.ReleaseToken(user, true)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	response.Success(ctx, gin.H{"accessToken": accessToken, "refreshToken": refreshToken}, "刷新成功")
}

func UserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": dto.ParseUserDTO(user.(model.User))}, "")
}

func isUsernameExist(db *gorm.DB, username string) bool {
	var user model.User
	db.Where("username = ?", username).First(&user)
	return user.ID != 0
}
