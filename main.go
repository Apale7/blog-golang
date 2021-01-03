package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	r := gin.Default()
	viper.SetConfigName("JWT_conf")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Error(errors.WithStack(err))
		panic("viper readInConfig error")
	}
	key := []byte(viper.GetString("key"))
	store := cookie.NewStore([]byte(key))
	r.Use(sessions.Sessions("login_session", store))
	r = CollectRouter(r)
	r.Run(":8081")
}

