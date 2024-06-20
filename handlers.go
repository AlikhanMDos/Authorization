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

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with a unique phone number
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User					true	"User details"
//	@Success		200		{object}	map[string]interface{}	"User registered successfully"
//	@Failure		400		{object}	map[string]string		"Invalid input"
//	@Failure		500		{object}	map[string]string		"Error creating user"
//	@Router			/register [post]
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

// Login godoc
//
//	@Summary		Login an existing user
//	@Description	Login an existing user with phone number and password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			user	body		User					true	"User credentials"
//	@Success		200		{object}	map[string]interface{}	"User logged in successfully"
//	@Failure		400		{object}	map[string]string		"Invalid input"
//	@Failure		401		{object}	map[string]string		"Invalid phone or password"
//	@Failure		500		{object}	map[string]string		"Error generating token"
//	@Router			/login [post]
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

//	@BasePath	/api/v1

// Logout godoc
//
//	@Summary		Logout the current user
//	@Description	Logout the current user by clearing the JWT token
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Logged out successfully"
//	@Router			/logout [post]
func Logout(c *gin.Context) {
	// Clear the JWT token from the client (for example, from cookies)
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// UpdatePassword godoc
//
//	@Summary		Update user password
//	@Description	Update the password for the authenticated user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"User ID"
//	@Param			user	body		User				true	"User details"
//	@Success		200		{object}	map[string]string	"Profile updated successfully"
//	@Failure		400		{object}	map[string]string	"Invalid ID or input"
//	@Failure		500		{object}	map[string]string	"Failed to update profile info"
//	@Router			/users/{id}/password [put]
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

// UserProfile godoc
//
//	@Summary		Get user profile
//	@Description	Get the profile of a user by ID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"User ID"
//	@Success		200	{object}	map[string]interface{}	"User profile"
//	@Failure		404	{object}	map[string]string		"User Id is required or Failed to fetch user profile"
//	@Router			/users/{id} [get]
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

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post for the authenticated user
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			post	body		Post				true	"Post details"
//	@Success		200		{object}	Post				"Post created successfully"
//	@Failure		400		{object}	map[string]string	"Invalid input"
//	@Failure		500		{object}	map[string]string	"Failed to create post"
//	@Router			/posts [post]
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

// GetPosts godoc
//
//	@Summary		Get all posts
//	@Description	Get all posts from the database
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		Post				"List of posts"
//	@Failure		500	{object}	map[string]string	"Failed to fetch posts or parse posts"
//	@Router			/posts [get]
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

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update an existing post by ID for the authenticated user
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Post ID"
//	@Param			post	body		Post				true	"Post details"
//	@Success		200		{object}	map[string]string	"Post updated successfully"
//	@Failure		400		{object}	map[string]string	"Invalid post ID or input"
//	@Failure		500		{object}	map[string]string	"Failed to update post"
//	@Router			/posts/{id} [put]
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

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete an existing post by ID for the authenticated user
//	@Tags			Post
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Post ID"
//	@Success		200	{object}	map[string]string	"Post deleted successfully"
//	@Failure		400	{object}	map[string]string	"Invalid post ID"
//	@Failure		500	{object}	map[string]string	"Failed to delete post"
//	@Router			/posts/{id} [delete]
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

// Protected godoc
//
//	@Summary		Access protected route
//	@Description	Access a protected route for authenticated users
//	@Tags			Example
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Welcome to the protected route!"
//	@Router			/protected [get]
func Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the protected route!"})
}
