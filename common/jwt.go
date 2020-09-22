package common

import (
	"blog-gin/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var accessKey, refreshKey []byte

func init() {
	viper.SetConfigName("JWT_conf")
	viper.AddConfigPath("./config")
	if err = viper.ReadInConfig(); err != nil {
		log.Error(errors.WithStack(err))
		panic("viper readInConfig error")
	}
	accessKey = []byte(viper.GetString("accessKey"))
	refreshKey = []byte(viper.GetString("refreshKey"))
	log.Println(string(accessKey), string(refreshKey))
}

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User, isRefresh bool) (string, error) {
	var expirationTime time.Time
	if isRefresh {
		expirationTime = time.Now().Add(120 * time.Minute)
	} else {
		expirationTime = time.Now().Add(30 * time.Minute)
	}

	//var claims *Claims
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "apale",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var tokenString string
	var err error
	if isRefresh {
		tokenString, err = token.SignedString(refreshKey)
	} else {
		tokenString, err = token.SignedString(accessKey)
	}
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string, isRefresh bool) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if isRefresh {
			return refreshKey, nil
		} else {
			return accessKey, nil
		}
	})
	return token, claims, err
}
