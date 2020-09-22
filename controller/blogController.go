package controller

import (
	"blog-gin/common"
	"blog-gin/model"
	"blog-gin/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func CreateBlog(ctx *gin.Context) {
	userInterface, _ := ctx.Get("user") //有中间件，所以一定存在
	user := userInterface.(model.User)
	blog := model.Blog{}
	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	err = json.Unmarshal(bytes, &blog)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求参数错误")
		return
	}
	if blog.Title == "" {
		response.Fail(ctx, nil, "标题不能为空")
		return
	}
	if len([]byte(blog.Title)) > 128 {
		response.Fail(ctx, nil, "标题过长")
		return
	}
	if blog.Content == "" {
		response.Fail(ctx, nil, "正文不能为空")
		return
	}
	blog.UserID = user.ID
	db := common.GetDb()
	db.Create(&blog)
	response.Success(ctx,nil,"新建博客成功")
}

func UpdateBlog(ctx *gin.Context) {
	userInterface, _ := ctx.Get("user") //有中间件，所以一定存在
	user := userInterface.(model.User)
	blog := model.Blog{}
	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	err = json.Unmarshal(bytes, &blog)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求参数错误")
		return
	}
	db := common.GetDb()
	oldBlog := model.Blog{}
	db.First(&oldBlog, blog.ID)
	if oldBlog.ID <= 0{
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求参数错误")
		return
	}
	if oldBlog.UserID != user.ID {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求参数错误")
		return
	}

	if blog.Title == "" {
		response.Fail(ctx, nil, "标题不能为空")
		return
	}
	if len([]byte(blog.Title)) > 128 {
		response.Fail(ctx, nil, "标题过长")
		return
	}
	if blog.Content == "" {
		response.Fail(ctx, nil, "正文不能为空")
		return
	}
	db.Model(&model.Blog{}).Updates(model.Blog{Title: blog.Title, Content: blog.Content})
	response.Success(ctx,nil,"修改博客成功")
}