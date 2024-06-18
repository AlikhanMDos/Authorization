package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword

	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser User
	err := userCollection.FindOne(context.TODO(), bson.M{"phone": user.Phone}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone or password"})
		return
	}

	if !CheckPasswordHash(user.Password, foundUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone or password"})
		return
	}
	token, err := GenerateJWT(user.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	// Clear the JWT token from the client (for example, from cookies)
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// UserProfile handles the GET /profile route
// UserProfile handles the GET /profile route
func UserProfile(c *gin.Context) {
	// Extract phone number from query parameter
	userId := c.Query("id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}

	// Find user by phone number in MongoDB
	var user User
	err := userCollection.FindOne(context.Background(), bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	// Respond with user's name and phone number
	c.JSON(http.StatusOK, gin.H{
		"name":  user.Name,
		"phone": user.Phone,
	})
}

// Extract token from Authorization header
//tokenString := extractTokenFromHeader(c.Request.Header.Get("Authorization"))
//if tokenString == "" {
//	c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
//	return
//}

//	// Add token to blacklist
//	mutex.Lock()
//	defer mutex.Unlock()
//	tokenBlacklist[tokenString] = struct{}{}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
//}

func Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected route!"})
}
