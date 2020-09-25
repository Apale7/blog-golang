package controller

import (
	"blog-gin/common"
	"blog-gin/dto"
	"blog-gin/model"
	"blog-gin/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
)

func CreateBlog(ctx *gin.Context) {
	user, blog, err := getUserAndBlog(ctx)
	if err != nil || !checkBlog(ctx, blog) {
		return
	}
	blog.UserID = user.ID
	blog.Content = html.EscapeString(blog.Content)
	blog.Title = html.EscapeString(blog.Title)
	db := common.GetDb()
	db.Create(&blog)
	response.Success(ctx, nil, "新建博客成功")
}

func UpdateBlog(ctx *gin.Context) {
	user, blog, err := getUserAndBlog(ctx)
	if err != nil || !checkUserId(ctx, blog, user) || checkBlog(ctx, blog) {
		return
	}
	blog.Content = html.EscapeString(blog.Content)
	blog.Title = html.EscapeString(blog.Title)
	db := common.GetDb()
	db.Model(&model.Blog{}).Updates(model.Blog{Title: blog.Title, Content: blog.Content})
	response.Success(ctx, nil, "修改博客成功")
}

func DeleteBlog(ctx *gin.Context) {
	user, blog, err := getUserAndBlog(ctx)
	if err != nil || !checkUserId(ctx, blog, user) {
		return
	}
	db := common.GetDb()
	db.Model(&model.Blog{}).Delete(&blog)
	response.Success(ctx, nil, "删除博客成功")
}

func getUserAndBlog(ctx *gin.Context) (model.User, model.Blog, error) { //从context获取user和blog
	userInterface, _ := ctx.Get("user") //有中间件，所以一定存在
	user := userInterface.(model.User)
	blog := model.Blog{} //这个blog的blogID可信，userID不可信
	bytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return model.User{}, model.Blog{}, err
	}
	err = json.Unmarshal(bytes, &blog)
	if err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "请求参数错误")
		return model.User{}, model.Blog{}, err
	}
	return user, blog, nil
}

func checkUserId(ctx *gin.Context, blog model.Blog, user model.User) bool { //判断是不是作者本人修改博客
	db := common.GetDb()
	oldBlog := model.Blog{}
	db.First(&oldBlog, blog.ID)
	if oldBlog.ID <= 0 {
		response.Fail(ctx, nil, "请求参数错误")
		return false
	}
	if oldBlog.UserID != user.ID {
		response.Fail(ctx, nil, "请求参数错误")
		return false
	}
	return true
}

func BlogList(ctx *gin.Context) { //返回一个用户的所有文章
	userid, err := strconv.ParseInt(ctx.DefaultQuery("userid", "0"), 10, 32)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	db := common.GetDb()
	user := model.User{}
	db.Model(&model.User{}).Where("id=?", userid).First(&user)
	if user.ID == 0 {
		response.Fail(ctx, nil, "请求参数错误")
		return
	}
	var blogs []model.Blog
	db.Model(&model.Blog{}).Select("id, title, create_at").Where("user_id=?", user.ID).Find(&blogs)
	blogDTOs := make([]dto.BlogDTO, len(blogs))
	for i := range blogs {
		blogDTOs[i] = dto.ParseBlogDTO(blogs[i])
	}
	response.Success(ctx, gin.H{"blogs": blogDTOs}, "查询博客列表成功")
}

func BlogContent(ctx *gin.Context) {
	blogId, err := strconv.ParseInt(ctx.DefaultQuery("blogid", "0"), 10, 32)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	db := common.GetDb()
	blog := model.Blog{}
	db.Model(&model.Blog{}).Select("id, title, created_at, content").First(&blog, blogId)
	if blog.ID > 0 {
		response.Success(ctx, gin.H{"blog": dto.ParseBlogDTO(blog)}, "查询博客列表成功")
	} else {
		response.Fail(ctx, nil, "请求参数错误")
	}
}

func checkBlog(ctx *gin.Context, blog model.Blog) bool {
	if blog.Title == "" {
		response.Fail(ctx, nil, "标题不能为空")
		return false
	}
	if len([]byte(blog.Title)) > 64 {
		response.Fail(ctx, nil, "标题过长")
		return false
	}
	if blog.Content == "" {
		response.Fail(ctx, nil, "正文不能为空")
		return false
	}
	return true
}
