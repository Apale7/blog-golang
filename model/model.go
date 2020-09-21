package model

import "github.com/jinzhu/gorm"

type Conf struct {
	Mysql Mysql `json:"mysql"`
}
type Mysql struct {
	Dbname   string `json:"dbname"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
}

type User struct {
	gorm.Model
	Username string `gorm:"size:32;unique;not null;index:username_idx"`
	Password string `gorm:"size:32;not null"`
}

type Blog struct {
	gorm.Model
	Title   string `gorm:"size:128;not null"`
	content string `gorm:"type:longtext;not null"`
	UserID uint
}

type Comment struct {
	gorm.Model
	Title   string `gorm:"size:128;not null"`
	content string `gorm:"size:256;not null"`
	BlogID  uint
	UserID uint
}