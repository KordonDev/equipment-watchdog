package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordondev/equipment-watchdog/security"
)

type member struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var members = []member{
	{Id: "1", Name: "Arne"},
	{Id: "2", Name: "Luka"},
}

func getMembers(c *gin.Context) {
	c.JSON(http.StatusOK, members)
}

func addMember(c *gin.Context) {
	var newMember member

	if err := c.BindJSON(&newMember); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fmt.Println(newMember)

	members = append(members, newMember)
	c.JSON(http.StatusCreated, newMember)
}

func main() {
	router := gin.Default()

	members := router.Group("/members", security.AuthorizeJWTMiddleware())
	members.GET("/", getMembers)
	members.POST("/", addMember)

	userDB := security.NewUserDB()
	webAuthNService := security.NewWebAuthNService(userDB)

	router.GET("/register/:username", webAuthNService.StartRegister)
	router.POST("register/:username", webAuthNService.FinishRegistration)

	router.GET("/login/:username", webAuthNService.StartLogin)
	router.POST("login/:username", webAuthNService.FinishLogin)

	router.LoadHTMLGlob("templates/*.html")

	router.GET("/index/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"name": name,
		})
	})

	router.Run("localhost:8080")
}
