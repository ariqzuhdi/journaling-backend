package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/cheeszy/go-crud/dto"
	"github.com/cheeszy/go-crud/initializers"
	"github.com/cheeszy/go-crud/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NotFoundHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"error": "URL not found.",
	})
}

func PostsCreate(c *gin.Context) {
	// Get data off requests body
	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID type"})
		return
	}

	post := models.Post{
		Title:  body.Title,
		Body:   body.Body,
		UserID: userID,
	}

	result := initializers.DB.Create(&post)
	if result.Error != nil {
		c.Status(400)
		return
	}

	initializers.DB.Preload("User").First(&post, post.ID)

	postResponse := dto.PostResponse{
		ID:    post.ID,
		Title: post.Title,
		Body:  post.Body,
	}

	// return it
	c.JSON(200, gin.H{
		"post": postResponse,
	})
}

func PostsShow(c *gin.Context) {

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Get the id of url
	id := c.Param("id")

	var post []models.Post
	result := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&post, id)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsShowAll(c *gin.Context) {

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Get the id of url
	username := c.Param("username")

	// take user from username
	var user models.User

	if err := initializers.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	if user.ID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Forbidden: You can't access this resource",
		})
		return
	}

	var posts []models.Post
	result := initializers.DB.Preload("User").Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}

	var postResponses []dto.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, dto.PostResponse{
			ID:    post.ID,
			Title: post.Title,
			Body:  post.Body,
			User:  dto.UserResponse{Username: post.User.Username},
		})
	}

	c.JSON(200, gin.H{
		"posts": postResponses,
	})
}

func PostsUpdate(c *gin.Context) {

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Get the id of url
	id := c.Param("id")

	// get the data off the req body

	var body struct {
		Body  string
		Title string
	}

	c.Bind(&body)

	// find the post were updating
	var post []models.Post
	result := initializers.DB.Where("id = ? AND user_id = ?", id, userID).First(&post, id)
	if result.Error != nil {
		c.JSON(401, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	initializers.DB.Where("id = ? AND user_id = ?", id, userID).Model(&post).Updates(models.Post{
		Title: body.Title, Body: body.Body,
	})

	// updating
	c.JSON(200, gin.H{
		"post": post,
	})
}

func PostsDelete(c *gin.Context) {

	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID type"})
		return
	}

	id := c.Param("id")
	var post models.Post

	result := initializers.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&post, id)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Failed to delete post"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Post not found or unauthorized"})
		return
	}

	c.JSON(200, gin.H{"message": "Deleted"})

}

func PostsIndex(c *gin.Context) {
	//Get the posts
	var posts []models.Post
	initializers.DB.Find(&posts)

	//Respond with them
	c.JSON(200, gin.H{
		"post": posts,
	})
}

func MonkeyAPI(c *gin.Context) {
	apiKey := os.Getenv("MONKEYTYPE_API_KEY")
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.monkeytype.com/users/personalBests?mode=time", nil)
	req.Header.Add("Authorization", "ApeKey "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read response body"})
		return
	}

	// Forward response eksternal ke client
	c.Data(resp.StatusCode, "application/json", body)
}
