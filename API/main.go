package main

import (
	svegoTypes "API/types"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	dotEnv "github.com/miguelmagueijo/golangDotEnv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*####################################################################################################################*/
/* Global variables */
/*####################################################################################################################*/
var DB *sql.DB
var DotEnvData map[string]string

/*####################################################################################################################*/
/* Start functions */
/*####################################################################################################################*/
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

/*####################################################################################################################*/
/* Helper functions */
/*####################################################################################################################*/
func extractJwtString(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")

	headerParts := strings.Split(bearerToken, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid authorization header")
	}

	return headerParts[1], nil
}

/*####################################################################################################################*/
/* Middlewares */
/*####################################################################################################################*/
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractJwtString(c.Request)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(DotEnvData["JWT_SIGNING_KEY"]), nil
		})

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		if !token.Valid {
			log.Println("Invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Could not get token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		expDate, err := claims.GetExpirationTime()
		if err != nil {
			log.Println("Could not get token expiration date")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		if expDate.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "expired session",
			})
			return
		}

		fmt.Println(claims)

		userId, ok := claims["user_id"].(string)
		if !ok {
			log.Println("Could not extract \"user_id\" from claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		c.Set("user_id", userId)
		c.Next()
	}
}

/*####################################################################################################################*/
/* Routers Handlers */
/*####################################################################################################################*/
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

func meHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data": map[string]interface{}{
			"user_id": c.MustGet("user_id"),
		},
	})
}

func loginHandler(c *gin.Context) {
	if _, err := extractJwtString(c.Request); err == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var formData svegoTypes.LoginFormData
	if err := c.Bind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	var userId int64
	var password string
	var isActive bool

	row := DB.QueryRow("SELECT id, password, is_active FROM user WHERE username = ?", formData.Username)

	err := row.Scan(&userId, &password, &isActive)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
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
			"error": "User not found",
		})
		return
	}

	tokenId, err := uuid.NewRandom()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong logging you in",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": strconv.FormatInt(userId, 10),
		"iss":     "svego.todolist.forgedin.space",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
		"jti":     tokenId.String(),
	})

	tokenString, err := token.SignedString([]byte(DotEnvData["JWT_SIGNING_KEY"]))

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong logging you in",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You have successfully logged in",
		"token":   tokenString,
	})
}

func main() {
	DotEnvData = dotEnv.Load()

	if len(DotEnvData["JWT_SIGNING_KEY"]) != 512 {
		log.Fatal("JWT signing key should be 512 characters long")
	}

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

	routerV1 := router.Group("/v1")

	routerV1.POST("/login", loginHandler)

	authorizedGroup := router.Group("/")
	authorizedGroup.Use(JwtAuthMiddleware())
	{
		v1Group := authorizedGroup.Group("v1")

		v1Group.GET("/tasks", getTasks)
		v1Group.GET("/me", meHandler)

		// nested group
		//testing := authorizedGroup.Group("testing")
		//testing.GET("/analytics", analyticsEndpoint)
	}

	err = router.Run(":4555")

	if err != nil {
		panic("[ERROR] Something went wrong while starting the server: " + err.Error())
	}
}
