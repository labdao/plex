package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "io/ioutil"
)

func main() {
  r := gin.Default()

  r.GET("/_health", health)
  r.POST("/judge", judge)

  r.Run()
}

func health(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"status": "ok"})    
}

func judge(c *gin.Context) {
  // for now, print the request body to std out so we can peep it
  body, _ := ioutil.ReadAll(c.Request.Body)
  println(string(body))
  // the judge endpoint always returns status 200 to accept all jobs (for now)
  c.JSON(http.StatusOK, gin.H{})    
}
