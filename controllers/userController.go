package controllers

import (
	"net/http"
	"os"
	"time"
	"userapp/helpers"
	"userapp/initializers"
	"userapp/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func UserIndex(c *gin.Context) {
	var users []models.User

	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to get users",
		})

		return
	}

	c.JSON(http.StatusOK, users)
}

func UserGet(c *gin.Context) {
	user, err := helpers.GetUserFromUri(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, user)
}

func UserSignUp(c *gin.Context) {
	var body struct {
		Username string
		Password string
		Email    string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash Password",
		})

		return
	}

	user := models.User{Username: body.Username, Email: body.Email, Password: string(hash)}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UserLogin(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	var user models.User

	initializers.DB.First(&user, "username = ?", body.Username)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username or password",
		})

		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username or password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to make token",
		})
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*3, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func UserEdit(c *gin.Context) {
	user, err := helpers.GetUserFromUri(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var body struct {
		Email string `form:"email"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	res := initializers.DB.Model(&user).Update("email", body.Email)

	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email is already used",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UserChangePassword(c *gin.Context) {
	user, errUser := helpers.GetUserFromUri(c)

	if errUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errUser.Error(),
		})

		return
	}

	var body struct {
		Password    string
		NewPassword string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if errPass != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})

		return
	}

	hash, errComp := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if errComp != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash Password",
		})

		return
	}

	initializers.DB.Model(&user).Update("password", string(hash))

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

func UserDelete(c *gin.Context) {
	user, err := helpers.GetUserFromUri(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	initializers.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User Deleted",
	})
}
