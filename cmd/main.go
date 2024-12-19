package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
)

var rdb *redis.Client
var ctx = context.Background()

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	r := gin.Default()

	r.GET("/health", checkHealth)
	r.GET("/attendance/:user_id", getAttendance)
	r.POST("/attendance/:user_id", markAttendance)

	r.Run(":8080")
}

func checkHealth(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func getAttendance(c *gin.Context) {
	userID := c.Param("user_id")

	count, err := rdb.Get(ctx, userID).Result()
	if err == redis.Nil {
		c.JSON(200, gin.H{"user_id": userID, "attendance_count": 0})
	} else if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при запросе к Redis"})
	} else {
		c.JSON(200, gin.H{"user_id": userID, "attendance_count": count})
	}
}

func markAttendance(c *gin.Context) {
	userID := c.Param("user_id")

	_, err := rdb.Incr(ctx, userID).Result()
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при записи посещения в Redis"})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Посещение для пользователя %s успешно зарегистрировано", userID)})
}
