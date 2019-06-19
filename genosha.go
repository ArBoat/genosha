package main

import (
  "github.com/gin-gonic/gin"
  "log"
  "net/http"
)

func main() {
  router := gin.Default()
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
  // However, this one will match /user/john/ and also /user/john/send
  // If no other routers match /user/john, it will redirect to /user/john/
  router.GET("/user/:name/*action", func(c *gin.Context) {
    name := c.Param("name")
    action := c.Param("action")
    //_ := c.DefaultQuery("firstname", "Guest")
    //_ = c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

    message := name + " is " + action
    c.String(http.StatusOK, message)
    aaa(message)
    log.Println("bbb", message)
  })

  // For each matched request Context will hold the route definition
  router.POST("/user/:name/*action", func(c *gin.Context) {
    //c.FullPath() == "/user/:name/*action" // true
  })
  router.Run() // listen and serve on 0.0.0.0:8080
}

func aaa(message string)  {
  message = message + message
  log.Println("aaa", message)
}