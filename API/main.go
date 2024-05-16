package main

import (
	svegoTypes "API/types"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

/* Global variables */
var DB *sql.DB

/* Routers Handlers */
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

func loginHandler(c *gin.Context) {
	var formData svegoTypes.LoginFormData

	err := c.Bind(&formData)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	var userId int64
	var password string
	var isActive bool

	row := DB.QueryRow("SELECT id, password, is_active FROM user WHERE username = ?", formData.Username)

	err = row.Scan(&userId, &password, &isActive)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user",
		})
		return
	case err == nil:
		break
	default:
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(formData.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user",
		})
		return
	}

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong login you in... try again later",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "everything looks fine",
		"token":   hex.EncodeToString(tokenBytes),
	})
}

func createDatabase() {
	fileBytes, err := os.ReadFile("./db/create.sql")

	if err != nil {
		log.Fatal(err)
	}

	sqlCommands := string(fileBytes)

	tx, err := DB.Begin()

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

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	DB, err = sql.Open("sqlite3", "./db/svego.sqlite")

	if err != nil {
		log.Fatal(err)
	}

	createDatabase()

	router := gin.Default()

	router.GET("/", status)
	router.GET("/status", status)
	router.GET("/v1/tasks", getTasks)
	router.POST("/v1/login", loginHandler)

	err = router.Run(":4555")

	if err != nil {
		panic("[ERROR] Something went wrong while starting the server: " + err.Error())
	}
}
