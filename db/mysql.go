package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Db *gorm.DB
var err error

const base string = "%s:%s@tcp(%s:%d)/%s?charset=utf8"

func init() {
	viper.SetConfigName("db_conf")
	viper.AddConfigPath("./config")
	if err = viper.ReadInConfig(); err != nil {
		log.Error(errors.WithStack(err))
		panic("viper readInConfig error")
	}
	var dbconf conf
	if err = viper.Unmarshal(&dbconf); err != nil {
		log.Error(errors.WithStack(err))
		panic("viper Unmarshal error")
	}
	mysql_conf := dbconf.Mysql
	Db, err = gorm.Open("mysql", fmt.Sprintf(base, mysql_conf.Username, mysql_conf.Password, mysql_conf.Host, mysql_conf.Port, mysql_conf.Dbname))
	if err != nil {
		log.Error(errors.WithStack(err))
		return
	} else {
		fmt.Println("database linked")
	}
	if !Db.HasTable(&User{}) {
		Db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&User{})
	}
	if !Db.HasTable(&Blog{}) {
		Db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Blog{})
		Db.Model(&Blog{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	}
	if !Db.HasTable(&Comment{}) {
		Db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Comment{})
		Db.Model(&Comment{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		Db.Model(&Comment{}).AddForeignKey("blog_id", "blogs(id)", "RESTRICT", "RESTRICT")
	}
	//Db, err := gorm.Open("mysql")
}

type conf struct {
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
