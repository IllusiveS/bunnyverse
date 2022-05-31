package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func bounce(c *gin.Context) {
	owner := c.Param("bunnyID")
	ownerId, _ := strconv.ParseUint(owner, 10, 64)

	var rabbit Rabbit
	DB.First(&rabbit, ownerId)
	rabbit.Bounces += 1
	DB.Save(&rabbit)

	c.JSON(http.StatusOK, rabbit)
}

func getCarrotsSinceLastTime(c *gin.Context) {
	owner := c.DefaultQuery("owner", "-1")
	rabbitOwner, _ := extractPlayerFromStringID(owner)

	var carrots []Carrot
	DB.Where("created_at > ", rabbitOwner.LastCarrotsCheck).Find(&carrots)

	updatePlayerCarrotTimeChecked(rabbitOwner)
	DB.Save(&rabbitOwner)

	c.JSON(http.StatusOK, carrots)
}

func getBunniesSinceLastTime(c *gin.Context) {
	//TODO move param to query
	ownerParam := c.Param("ownerID")

	rabbitOwner, err := extractPlayerFromStringID(ownerParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var rabbits []Rabbit
	DB.Where("created_at > ?", rabbitOwner.LastBunniesCheck).Find(&rabbits)

	updatePlayerEnergy(rabbitOwner)
	rabbitOwner.LastBunniesCheck = time.Now()

	DB.Save(&rabbitOwner)

	c.JSON(http.StatusOK, rabbits)
}

func procreate(c *gin.Context) {

}
