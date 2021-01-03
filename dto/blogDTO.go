package dto

import (
	"blog-gin/common"
	"blog-gin/model"
)

type BlogDTO struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	Nickname  string `json:"nickname"`
}

func ParseBlogDTO(blog model.Blog) *BlogDTO {
	ret := &BlogDTO{
		Id:        blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		CreatedAt: blog.CreatedAt.Unix(),
	}
	db := common.GetDb()
	user := model.User{}
	db.Model(&model.User{}).Select("nickname").Where("id=?", blog.UserID).First(&user)
	ret.Nickname = user.Nickname
	return ret
}
