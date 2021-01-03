package controller

import (
	"blog-gin/common"
	"blog-gin/dto"
	"blog-gin/model"
	"blog-gin/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"html"
	"net/http"
	"strconv"
)

func CreateBlog(c *gin.Context) {
	user, blog, err := getUserAndBlog(c)
	if err != nil || !checkBlog(c, blog) {
		return
	}
	log.Printf("%+v", blog)
	blog.UserID = user.ID
	blog.Content = html.EscapeString(blog.Content)
	blog.Title = html.EscapeString(blog.Title)
	db := common.GetDb()
	db.Create(&blog)
	response.Success(c, nil, "新建博客成功")
}

func UpdateBlog(c *gin.Context) {
	user, blog, err := getUserAndBlog(c)
	if err != nil || !checkUserId(c, blog, user.ID) || checkBlog(c, blog) {
		return
	}
	blog.Content = html.EscapeString(blog.Content)
	blog.Title = html.EscapeString(blog.Title)
	db := common.GetDb()
	db.Model(&model.Blog{}).Updates(model.Blog{Title: blog.Title, Content: blog.Content})
	response.Success(c, nil, "修改博客成功")
}

func DeleteBlog(c *gin.Context) {
	user, blog, err := getUserAndBlog(c)
	if err != nil || !checkUserId(c, blog, user.ID) {
		return
	}
	db := common.GetDb()
	db.Model(&model.Blog{}).Delete(&blog)
	response.Success(c, nil, "删除博客成功")
}

func getUserAndBlog(c *gin.Context) (*model.User, *model.Blog, error) { //从context获取user和blog
	userInterface, _ := c.Get("user") //有中间件，所以一定存在
	user := userInterface.(model.User)
	blog := model.Blog{} //这个blog的blogID可信，userID不可信
	err := c.Bind(&blog)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "请求参数错误")
		return nil, nil, err
	}
	fmt.Printf("%+v", blog)
	return &user, &blog, nil
}

func checkUserId(c *gin.Context, blog *model.Blog, userID uint) bool { //判断是不是作者本人修改博客
	db := common.GetDb()
	oldBlog := model.Blog{}
	db.First(&oldBlog, blog.ID)
	if oldBlog.ID <= 0 {
		response.Fail(c, nil, "请求参数错误")
		return false
	}
	if oldBlog.UserID != userID {
		response.Fail(c, nil, "请求参数错误")
		return false
	}
	return true
}

func BlogList(c *gin.Context) { //返回一个用户的所有文章
	userid, err := strconv.ParseInt(c.DefaultQuery("userid", "0"), 10, 32)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	db := common.GetDb()
	var blogs []model.Blog

	if userid == 0 {
		db.Model(&model.Blog{}).Select("id, title, created_at, content,user_id").Order("created_at desc").Find(&blogs)
	} else {
		user := model.User{}
		db.Model(&model.User{}).Where("id=?", userid).First(&user)
		if user.ID == 0 {
			response.Fail(c, nil, "请求参数错误")
			return
		}
		db.Model(&model.Blog{}).Select("id, title, created_at, content,user_id").Where("user_id=?", user.ID).Find(&blogs)
	}
	blogDTOs := make([]*dto.BlogDTO, len(blogs))
	for i := range blogs {
		blogDTOs[i] = dto.ParseBlogDTO(blogs[i])
	}
	response.Success(c, gin.H{"blogs": blogDTOs}, "查询博客列表成功")
}

func BlogContent(c *gin.Context) {
	blogId, err := strconv.ParseInt(c.DefaultQuery("id", "0"), 10, 32)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Error(errors.WithStack(err))
		return
	}
	db := common.GetDb()
	blog := model.Blog{}
	db.Model(&model.Blog{}).Select("id, title, created_at, content").First(&blog, blogId)
	if blog.ID > 0 {
		response.Success(c, gin.H{"blog": dto.ParseBlogDTO(blog)}, "查询博客列表成功")
	} else {
		response.Fail(c, nil, "请求参数错误")
	}
}

func checkBlog(c *gin.Context, blog *model.Blog) bool {
	if blog.Title == "" {
		response.Fail(c, nil, "标题不能为空")
		return false
	}
	if len([]byte(blog.Title)) > 64 {
		response.Fail(c, nil, "标题过长")
		return false
	}
	if blog.Content == "" {
		response.Fail(c, nil, "正文不能为空")
		return false
	}
	return true
}
