package dto

import (
	"blog-gin/model"
	"time"
)

type BlogDTO struct {
	Id        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserId    uint      `json:"user_id"`
}

func ParseBlogDTO(blog model.Blog) BlogDTO {
	return BlogDTO{
		Id:        blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		CreatedAt: blog.CreatedAt,
		UserId:    blog.UserID,
	}
}
