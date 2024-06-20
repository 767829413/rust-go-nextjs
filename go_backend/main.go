package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required"`
}

func main() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("USER"),
		os.Getenv("DBNAME"),
		os.Getenv("PASSWORD"),
	)
	// dsn := fmt.Sprintf(
	// 	"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=Asia/Shanghai",
	// 	"localhost",
	// 	"5432",
	// 	"postgres",
	// 	"postgres",
	// 	"postgres",
	// )
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/api/rust/users", func(c *gin.Context) {
		var user User
		// 将请求体中的 JSON 数据解析到 user 结构体中
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "create user success",
		})
	})
	r.PUT("/api/rust/users/:id", func(c *gin.Context) {
		// 获取路由参数
		id := c.Param("id")
		var user User
		// 将请求体中的 JSON 数据解析到 user 结构体中
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "update user success",
		})
	})
	r.GET("/api/rust/users", func(c *gin.Context) {
		var users = []User{}
		if err := db.Raw("SELECT id, name, email FROM users").Find(&users).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})
	r.GET("/api/rust/users/:id", func(c *gin.Context) {
		var user = User{}
		// 获取路由参数
		id := c.Param("id")
		if err := db.Raw("SELECT id, name, email FROM users WHERE id = ?", id).Find(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	})
	r.DELETE("/api/rust/users/:id", func(c *gin.Context) {
		// 获取路由参数
		id := c.Param("id")
		if err := db.Exec("DELETE FROM users WHERE id = ?", id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "delete user success",
		})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
