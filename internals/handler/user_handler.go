package handler

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"filmfolk/internals/db"
	"filmfolk/internals/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims defines the structure of the JWT claims.
type Claims struct {
	UserID uint64          `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}




type RegisterRequest struct{
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// User Register Function

func RegisterUser(c *gin.Context) {

	var req RegisterRequest

	//converting JSON to Struct
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: "+ err.Error() })
		return
	}

	//Hashing Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error": "Failed to hash Password: "+ err.Error() })
		return
	}

	user := models.User{
		UserName:     req.UserName,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleUser,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + result.Error.Error()})
		return
	}


	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})



}

//Login Logic

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}




//login function

func LoginUser(c *gin.Context){

var req LoginRequest

if err := c.ShouldBindJSON(&req); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
}

var user models.User

result := db.DB.Where("email = ?", req.Email).First(&user)
if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash),[]byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

// JWT TOKEN

expirationTime := time.Now().Add(24 * time.Hour)

claims := &Claims{
	UserID: user.ID,
	Role:   user.Role,
	RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   strconv.FormatUint(user.ID, 10),
	},
}

jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
if len(jwtKey) == 0 {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret key is not configured on server"})
	return
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString(jwtKey)
if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}


	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}