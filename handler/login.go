package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"DOROserver/db"
	"DOROserver/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 註冊處理器
func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤"})
		return
	}

	collection := db.MongoClient.Database("doro").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 確認帳號是否存在
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Err()
	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusConflict, gin.H{"error": "帳號已存在"})
		return
	}

	// 產生 token 並寫入資料庫
	user.Token = GenerateToken()
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "註冊失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "註冊成功"})
}

// 登入處理器
func Login(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤"})
		return
	}

	collection := db.MongoClient.Database("doro").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		return
	}
	// 登入成功後設定 Cookie 儲存 token
	c.SetCookie("token", user.Token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "登入成功"})
}

// 登出處理器
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "已登出"})
}

// 產生隨機 token 字串
func GenerateToken() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "defaulttoken"
	}
	return hex.EncodeToString(bytes)
}
