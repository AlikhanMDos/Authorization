package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

// @BasePath /api/v1

// Register PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
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

	c.JSON(http.StatusOK, gin.H{
		"message":  "User registered successfully",
		"userData": user,
	})
}

// @BasePath /api/v1

// Login PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
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

	c.JSON(http.StatusOK, gin.H{"token": token, "userData": user})
}

// @BasePath /api/v1

// Logout PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Logout(c *gin.Context) {
	// Clear the JWT token from the client (for example, from cookies)
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func UpdatePassword(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the profile belongs to the authenticated user
	authorID := c.MustGet("userID").(primitive.ObjectID)
	filter := bson.M{"_id": objID, "author_id": authorID}

	update := bson.M{
		"$set": bson.M{
			"name":     user.Name,
			"phone":    user.Phone,
			"password": user.Password,
		},
	}

	result := userCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// @BasePath /api/v1

// UserProfile PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func UserProfile(c *gin.Context) {
	// Extract phone number from query parameter
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User Id is required"})
		return
	}

	// Find user by phone number in MongoDB
	var user User
	err1 := userCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err1 != nil {
		log.Println("Error fetching user profile:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch user profile"})
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

func CreatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.ID = primitive.NewObjectID()
	post.Date = time.Now()
	post.AuthorID = c.MustGet("userID").(primitive.ObjectID) // Assuming you have userID in context

	_, err := postCollection.InsertOne(context.Background(), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func GetPosts(c *gin.Context) {
	cursor, err := postCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer cursor.Close(context.Background())

	var posts []Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the post belongs to the authenticated user
	authorID := c.MustGet("userID").(primitive.ObjectID)
	filter := bson.M{"_id": objID, "author_id": authorID}

	update := bson.M{
		"$set": bson.M{
			"title": post.Title,
			"text":  post.Text,
			"image": post.Image,
			"date":  time.Now(),
		},
	}

	result := postCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

func DeletePost(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Ensure the post belongs to the authenticated user
	authorID := c.MustGet("userID").(primitive.ObjectID)
	filter := bson.M{"_id": objID, "author_id": authorID}

	result := postCollection.FindOneAndDelete(context.Background(), filter)
	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected route!"})
}
