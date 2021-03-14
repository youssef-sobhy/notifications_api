package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/youssef1337/notifications_api/tree/main/internal"
)

func main() {
	r := gin.Default()
	r.POST("api/v1/notifications/send", func(c *gin.Context) {
		body := internal.SendParams{}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ids, err := internal.CreateNotifications(body.Data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(200, ids)
	})

	r.GET("api/v1/notifications", func(c *gin.Context) {
		notifications, err := internal.FetchNotifications()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, notifications)
	})

	r.Run(":3000")
}
