package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	r := gin.Default()
	r.POST("api/v1/notifications/send", func(c *gin.Context) {
		body := SendParams{}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ids, err = CreateNotifications(body.Data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(200, ids)
	})

	r.GET("api/v1/notifications", func(c *gin.Context) {
		notifications, err := FetchNotifications()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, notifications)
	})

	r.Run(":3000")
}
