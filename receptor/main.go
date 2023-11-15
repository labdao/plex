package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labdao/receptor/models"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
)

func main() {
	log.Print("Connecting to database")
	models.ConnectDatabase()

	r := gin.Default()

	log.Print("Setting up routes")
	r.GET("/_health_check", health)
	r.POST("/judge", judge)

	r.Run()
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func judge(c *gin.Context) {

	// createa job row
	job := models.Job{}

	// normally want to catch an err here, but for now we always want to return 200
	c.BindJSON(&job)

	models.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&job)

	// the judge endpoint always returns status 200 to accept all jobs (for now)
	c.JSON(http.StatusOK, gin.H{})
}
