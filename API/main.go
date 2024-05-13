package main

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
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

func createDatabase(db *sql.DB) {
	fileBytes, err := os.ReadFile("./db/create.sql")

	if err != nil {
		log.Fatal(err)
	}

	sqlCommands := string(fileBytes)

	tx, err := db.Begin()

	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(sqlCommands)
	if err != nil {
		rollErr := tx.Rollback()

		if rollErr != nil {
			log.Fatal(rollErr)
		}

		log.Fatal(err)
	}

	err = tx.Commit()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := os.Remove("./db/svego.sqlite")

	if !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", "./db/svego.sqlite")

	if err != nil {
		log.Fatal(err)
	}

	createDatabase(db)

	defer db.Close()

	router := gin.Default()

	router.GET("/", status)
	router.GET("/status", status)
	router.GET("/tasks", getTasks)

	err = router.Run(":4555")

	if err != nil {
		panic("[ERROR] Something went wrong while starting the server: " + err.Error())
	}
}
