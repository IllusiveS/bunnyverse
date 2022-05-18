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
	r.GET("/myBunnies/:ownerID", getBunniesByOwner)
	r.POST("/procreate/:bunny1/:bunny2", procreate)

	r.Run()
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
