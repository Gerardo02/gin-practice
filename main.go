package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserID string

type Users struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var userData = make(map[UserID]Users)

func main() {
	r := gin.Default()

	r.GET("/ping", pong)

	r.GET("/users", getUsers)
	r.POST("/users", addUser)

	r.GET("/users/:id", getUserByID)
	r.PUT("/users/:id", editUser)
	r.DELETE("/users/:id", deleteUser)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func getUsers(c *gin.Context) {
	if len(userData) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "no users found",
		})
		return
	}

	users := []Users{}
	for _, user := range userData {
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "succesfull fetch",
		"data":    users,
		"length":  len(users),
	})
}

func editUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no id sent",
		})
		return
	}

	_, ok := userData[UserID(id)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}
	newUser := Users{}

	err := json.NewDecoder(c.Request.Body).Decode(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error creating user",
		})
		return
	}

	newUser.ID = id
	userData[UserID(id)] = newUser

	c.JSON(http.StatusOK, gin.H{
		"message": "user edited",
	})
}

func addUser(c *gin.Context) {
	newUser := Users{}

	err := json.NewDecoder(c.Request.Body).Decode(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error creating user",
		})
		return
	}

	userData[UserID(newUser.ID)] = newUser

	c.JSON(http.StatusOK, gin.H{
		"message": "user created",
	})
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no id sent",
		})
		return
	}

	user, ok := userData[UserID(id)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "succesfull",
		"data":    user,
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no id sent",
		})
		return
	}

	delete(userData, UserID(id))
	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted",
	})
}
