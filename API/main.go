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
	"time"
)

/*####################################################################################################################*/
/* Global variables */
/*####################################################################################################################*/
var DB *sql.DB
var DotEnvData map[string]string
var AccessTokenCookieDuration = 60 * 60
var RefreshTokenCookieDuration = getDaysInSeconds(30)

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
/* Database functions */
/*####################################################################################################################*/
func addRefreshToken(userId int64, jwtId string) (string, error) {
	refreshToken, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	stmt, err := DB.Prepare("INSERT INTO refresh_token (id, jwt_id, user_id) VALUES (?, ?, ?);")

	if err != nil {
		return "", err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(refreshToken.String(), jwtId, userId); err != nil {
		return "", err
	}

	return refreshToken.String(), nil
}

func setIsActiveRefreshToken(id string, isActive bool) error {
	tx, err := DB.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("UPDATE refresh_token SET is_active = ?, updated_at = unixepoch('now') WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(isActive, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errRollback
		}

		return err
	}

	return nil
}

func getRefreshTokenData(id string) (*svegoTypes.RefreshTokenData, error) {
	row := DB.QueryRow("SELECT id, jwt_id, user_id, is_active, expires_at, created_at, updated_at FROM refresh_token WHERE id = ?", id)

	tokenData := svegoTypes.RefreshTokenData{}
	err := row.Scan(&tokenData.Id, &tokenData.JwtId, &tokenData.UserId, &tokenData.IsActive, &tokenData.ExpiresAt, &tokenData.CreatedAt, &tokenData.UpdatedAt)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, err
	case err == nil:
		break
	default:
		log.Fatal(err)
	}

	return &tokenData, nil
}

func getUserJwtData(id int64) (*svegoTypes.UserJwtData, error) {
	row := DB.QueryRow("SELECT id, name, username, email, is_active FROM user WHERE id = ?;", id)

	userData := svegoTypes.UserJwtData{}
	var isActive bool
	err := row.Scan(&userData.Id, &userData.Name, &userData.Username, &userData.Email, &isActive)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, err
	case err == nil:
		break
	default:
		log.Fatal(err)
	}

	if !isActive {
		return nil, errors.New("cannot fetch jwt data from an user that is disabled")
	}

	return &userData, nil
}

func getUserMeData(id int64) (*svegoTypes.UserMeData, error) {
	row := DB.QueryRow("SELECT id, name, username, email, created_at, updated_at FROM user WHERE id = ?;", id)

	userData := svegoTypes.UserMeData{}
	err := row.Scan(&userData.Id, &userData.Name, &userData.Username, &userData.Email, &userData.CreatedAt, &userData.UpdatedAt)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, err
	case err == nil:
		break
	default:
		log.Fatal(err)
	}

	return &userData, nil
}

/*####################################################################################################################*/
/* Auxiliary functions */
/*####################################################################################################################*/
func generateJwt(userData *svegoTypes.UserJwtData, tokenId string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id":  strconv.FormatInt(userData.Id, 10),
		"username": userData.Username,
		"email":    userData.Email,
		"iat":      now.Unix(),
		"exp":      now.Add(time.Hour * 24 * 30).Unix(),
		"jti":      tokenId,
	})

	return token.SignedString([]byte(DotEnvData["JWT_SIGNING_KEY"]))
}

func getDaysInSeconds(n int) int {
	return 60 * 60 * 24 * n
}

func parseJwt(jwtString string) (*jwt.Token, error) {
	return jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return []byte(DotEnvData["JWT_SIGNING_KEY"]), nil
	})
}

/*####################################################################################################################*/
/* Request changer functions */
/*####################################################################################################################*/
func renewJwtAndInjectCookie(c *gin.Context, accessToken *string) error {
	refreshToken, err := c.Cookie("refreshToken")

	if err != nil {
		return err
	}

	refreshTokenData, err := getRefreshTokenData(refreshToken)

	if err != nil {
		return err
	}

	if !refreshTokenData.IsActive {
		return errors.New("refresh token expired")
	}

	if time.Now().After(time.Unix(refreshTokenData.ExpiresAt, 0)) {
		if err = setIsActiveRefreshToken(refreshToken, false); err != nil {
			return errors.New("refresh token expired")
		}
	}

	userData, err := getUserJwtData(refreshTokenData.UserId)

	if err != nil {
		return err
	}

	newJwt, err := generateJwt(userData, refreshTokenData.JwtId)

	if err != nil {
		return err
	}

	*accessToken = newJwt
	c.SetCookie("token", newJwt, AccessTokenCookieDuration, "/", "localhost", false, true)

	return nil
}

/*####################################################################################################################*/
/* Middlewares */
/*####################################################################################################################*/
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("token")

		if err != nil {
			if errRenew := renewJwtAndInjectCookie(c, &accessToken); errRenew != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				return
			}
		}

		token, err := parseJwt(accessToken)

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		if !token.Valid {
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
			if errRenew := renewJwtAndInjectCookie(c, &accessToken); errRenew != nil {
				log.Println(errRenew)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "couldn't renew token",
				})
				return
			} else {
				c.Redirect(http.StatusFound, "/logout")
				c.Abort()
				return
			}
		}

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
	fmt.Println("Getting tasks for user:", c.MustGet("user_id"))

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "not yet implemented",
	})
}

func meHandler(c *gin.Context) {
	userIdStr := c.MustGet("user_id").(string)

	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not fetch user data",
		})
		return
	}

	userData, err := getUserMeData(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not fetch user data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    userData,
	})
}

func logoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "localhost", false, true)

	c.Redirect(http.StatusFound, "/")
}

func loginHandler(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	if _, err := c.Cookie("refreshToken"); err == nil {
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

	var userJwtData svegoTypes.UserJwtData
	var password string
	var isActive bool
	row := DB.QueryRow("SELECT id, username, email, password, is_active FROM user WHERE username = ?", formData.Username)

	err := row.Scan(&userJwtData.Id, &userJwtData.Username, &userJwtData.Email, &password, &isActive)

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

	if !isActive {
		log.Println("Disabled user trying to login.", userJwtData)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(formData.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User not found",
		})
		return
	}

	jwtId, err := uuid.NewRandom()

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong logging you in",
		})
		return
	}

	refreshToken, err := addRefreshToken(userJwtData.Id, jwtId.String())

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong logging you in",
		})
		return
	}

	jwtString, err := generateJwt(&userJwtData, jwtId.String())

	c.SetCookie("token", jwtString, AccessTokenCookieDuration, "/", "localhost", false, true)
	c.SetCookie("refreshToken", refreshToken, RefreshTokenCookieDuration, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "You have successfully logged in",
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
	routerV1.GET("/logout", logoutHandler)

	authorizedGroup := router.Group("/")
	authorizedGroup.Use(JwtAuthMiddleware())
	{
		v1Group := authorizedGroup.Group("v1")

		v1Group.GET("/tasks", getTasks)
		v1Group.GET("/me", meHandler)
	}

	err = router.Run(":4555")

	if err != nil {
		panic("[ERROR] Something went wrong while starting the server: " + err.Error())
	}
}
