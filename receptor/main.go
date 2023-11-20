package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labdao/receptor/models"
	"gorm.io/gorm/clause"
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
	var requestPayload models.JobModel

	if err := c.BindJSON(&requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract Job.ID from the JSON data
	var jobID struct{ ID string }
	if err := json.Unmarshal(requestPayload.Spec, &jobID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new JobModel instance
	jobModel := models.JobModel{
		NodeID: requestPayload.NodeID,
		Spec:   requestPayload.Spec,
		JobID:  jobID.ID,
	}

	// Create or update the record in the database
	models.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&jobModel)

	c.JSON(http.StatusOK, gin.H{})
}
