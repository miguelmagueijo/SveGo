package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func getTasks(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "not yet implemented",
	})
}

func main() {
	router := gin.Default()

	router.GET("/", status)
	router.GET("/status", status)
	router.GET("/tasks", getTasks)

	err := router.Run(":4555")

	if err != nil {
		panic("[ERROR] Something went wrong while starting the server: " + err.Error())
	}
}
