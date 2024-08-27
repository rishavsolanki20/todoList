package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Topic struct {
	ID          uint   `gorm:"primaryKey"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
}

func main() {
	dsn := "user=postgres.vlhilrmgqjuxaibhwave password=!@#Rishav#@! host=aws-0-ap-southeast-1.pooler.supabase.com port=6543 dbname=postgres"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&Topic{})

	r := gin.Default()

	r.GET("/view", func(c *gin.Context) {
		var topic []Topic
		db.Find(&topic)
		c.JSON(http.StatusOK, topic)
	})

	r.POST("/create", func(c *gin.Context) {
		var newTopic Topic
		if err := c.ShouldBindJSON(&newTopic); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&newTopic).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, newTopic)
	})

	r.PUT("/update/:id", func(c *gin.Context) {
		id := c.Param("id")

		var topic Topic
		if err := db.First(&topic, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Topic not Found"})
		}

		if err := c.ShouldBindJSON(&topic); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Save(&topic).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, topic)
	})

	r.DELETE("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Delete(&Topic{}, id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "Topic deleted successfully")
	})
	r.Run(":8080")
}
