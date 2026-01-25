package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Print("Hello test")

	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "go url live air test test test",
		})
	})

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("failed to start the web server - Errror: %v", err))
	}
}
