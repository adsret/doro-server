package main

import (
	"log"
	"net/http"

	"DOROserver/db"
	"DOROserver/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	if err := db.ConnectMongo(); err != nil {
		log.Fatal("❌ Mongo 連線失敗:", err)
	}

	// 正確的 HTML 載入目錄
	r.LoadHTMLGlob("templates/html/*.html")
	r.Static("/img", "./templates/img")
	r.Static("/css", "./templates/css")

	// 功能路由
	r.POST("/upload", handler.UploadImage)
	r.GET("/myroom", func(c *gin.Context) { c.HTML(http.StatusOK, "myroom.html", nil) })
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.GET("/logout", handler.Logout)

	// 首頁
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

	r.Run("0.0.0.0:55688")
}
