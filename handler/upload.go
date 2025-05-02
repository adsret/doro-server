package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"DOROserver/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func UploadImage(c *gin.Context) {
	// 從 cookie 取得 token
	token, err := c.Cookie("token")
	if err != nil || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登入或缺少 token"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCol := db.MongoClient.Database("doro").Collection("users")

	var user bson.M
	err = userCol.FindOne(ctx, bson.M{"token": token}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 無效"})
		return
	}

	username := user["username"].(string)
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "上傳失敗: " + err.Error()})
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "讀取檔案失敗"})
		return
	}

	imgCol := db.MongoClient.Database("doro").Collection("user_images")
	_, err = imgCol.InsertOne(ctx, bson.M{
		"username": username,
		"filename": header.Filename,
		"filetype": header.Header.Get("Content-Type"),
		"data":     buf.Bytes(),
		"uploaded": time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "儲存資料庫失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "圖片上傳成功 ✅"})
}
