package common

import (
	"blog-gin/model"
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
	var dbconf model.Conf
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
	if !Db.HasTable(&model.User{}) {
		Db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&model.User{})
	}
	if !Db.HasTable(&model.Blog{}) {
		Db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&model.Blog{})
		Db.Model(&model.Blog{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	}
	if !Db.HasTable(&model.Comment{}) {
		Db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&model.Comment{})
		Db.Model(&model.Comment{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
		Db.Model(&model.Comment{}).AddForeignKey("blog_id", "blogs(id)", "RESTRICT", "RESTRICT")
	}
	//Db, err := gorm.Open("mysql")
}

func GetDb() *gorm.DB {
	return Db
}

