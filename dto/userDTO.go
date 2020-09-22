package dto

import "blog-gin/model"

type UserDTO struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func ParseUserDTO(user model.User) UserDTO {
	return UserDTO{
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}
}
