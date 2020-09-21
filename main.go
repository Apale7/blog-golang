package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r = CollectRouter(r)
	r.Run(":8080")
}

