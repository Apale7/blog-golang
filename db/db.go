package db

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type mysqlConf struct {
	Dbname string
	Host string
	Port int
	Password string
}

func init() {
	viper.SetConfigName("db_conf")
	viper.AddConfigPath("./config")
	if err:=viper.ReadInConfig(); err!=nil {
		log.Error(errors.WithStack(err))
	}
	var dbconf mysqlConf
	if err := mapstructure.Decode(viper.Get("mysql"), &dbconf); err != nil {
		log.Error(errors.WithStack(err))
	}
	fmt.Println(dbconf.Host)
	//Db, err := gorm.Open("mysql")
}