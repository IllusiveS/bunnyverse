package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	ConnectDatabase()

	r := gin.Default()
	r.Use()

	constantsGroup := r.Group("/constants")
	constantsGroup.GET("/energyRechargeRate", getEnergyRechargeRate)
	constantsGroup.GET("/carrotEnergyCost", getCarrotEnergyCost)
	constantsGroup.GET("/carrotLifetime", getCarrotLifetime)
	constantsGroup.GET("/bunnyEnergyCost", getBunnyEnergyCost)

	getNewGroup := r.Group("/getNew")
	getNewGroup.GET("/bunnies/:ownerID", getBunniesSinceLastTime)
	getNewGroup.GET("/carrots/", getCarrotsSinceLastTime)

	creationGroup := r.Group("/create")
	creationGroup.POST("/bunny/:bunnyName", addBunny)
	creationGroup.POST("/carrot", addCarrot)
	creationGroup.POST("/user", connectUser)
	creationGroup.POST("/bounce/:bunnyID", bounce)

	r.GET("/", getBunnies)
	r.GET("/energy", getEnergy)
	r.GET("/myBunnies/:ownerID", getBunniesByOwner)

	r.POST("/procreate/:bunny1/:bunny2", procreate)

	r.Run()
}

func getEnergy(c *gin.Context) {
	player, err := extractPlayerFromIDContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var rabbitsOwned []Rabbit
	DB.Where("Owner = ?", player.ID).Find(&rabbitsOwned)

	var carrotsCreated []Carrot
	DB.Where("Owner = ?", player.ID).Find(&carrotsCreated)

	updatePlayerEnergy(player)

	energy := player.Energy

	for _, rabbit := range rabbitsOwned {
		energy -= rabbit.EnergyCost
	}

	for _, carrot := range carrotsCreated {
		energy -= carrot.EnergyCost
	}

	DB.Save(player)

	c.JSON(http.StatusOK, gin.H{
		"energy": energy,
	})
}

func getBunniesByOwner(c *gin.Context) {
	owner := c.Param("ownerID")

	var rabbitsOwned []Rabbit
	DB.Where("Owner = ?", owner).Find(&rabbitsOwned)

	c.JSON(http.StatusOK, rabbitsOwned)
}

func getBunnies(c *gin.Context) {
	var rabbits []Rabbit
	DB.Find(&rabbits)

	c.JSON(http.StatusOK, rabbits)
}
