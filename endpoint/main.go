package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var energyRechargeRate = 1.0
var carrotEnergyCost = 10.0
var carrotLifetime = time.Second * 10
var bunnyEnergyCost = 100.0

func main() {
	ConnectDatabase()

	r := gin.Default()
	r.Use()

	constantsGroup := r.Group("/constants")
	constantsGroup.GET("/energyRechargeRate", getEnergyRechargeRate)
	constantsGroup.GET("/carrotEnergyCost", getCarrotEnergyCost)
	constantsGroup.GET("/carrotLifetime", getCarrotLifetime)
	constantsGroup.GET("/bunnyEnergyCost", getBunnyEnergyCost)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", getBunnies)
	r.GET("/newBunnies/:ownerID", getBunniesSinceLastTime)
	r.GET("/myBunnies/:ownerID", getBunniesByOwner)
	//r.POST("/bounce")
	r.POST("/new/:bunnyName", addBunny)
	r.POST("/new/carrot", addCarrot)
	r.POST("/bounce/:bunnyID", bounce)
	r.POST("/procreate/:bunny1/:bunny2", procreate)

	r.Run()
}

func procreate(c *gin.Context) {

}

func getEnergyRechargeRate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"energyRechargeRate": energyRechargeRate,
	})
}
func getCarrotEnergyCost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"carrotEnergyCost": carrotEnergyCost,
	})
}
func getCarrotLifetime(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"carrotLifetime": carrotLifetime,
	})
}
func getBunnyEnergyCost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"bunnyEnergyCost": bunnyEnergyCost,
	})
}

func bounce(c *gin.Context) {
	owner := c.Param("bunnyID")
	ownerId, _ := strconv.ParseUint(owner, 10, 64)

	var rabbit Rabbit
	DB.First(&rabbit, ownerId)
	rabbit.Bounces += 1
	DB.Save(&rabbit)

	c.JSON(http.StatusOK, rabbit)
}

func getBunniesByOwner(c *gin.Context) {
	owner := c.Param("ownerID")

	var rabbitsOwned []Rabbit
	DB.Where("Owner = ?", owner).Find(&rabbitsOwned)

	c.JSON(http.StatusOK, rabbitsOwned)
}

func extractCoordinates(c *gin.Context) (float64, float64, error) {
	xCoordString := c.Query("xCoord")
	yCoordString := c.Query("yCoord")

	xCoord, err := strconv.ParseFloat(xCoordString, 10)
	if err != nil {
		return 0, 0, err
	}

	yCoord, xErr := strconv.ParseFloat(yCoordString, 10)
	if xErr != nil {
		return 0, 0, xErr
	}

	return xCoord, yCoord, nil
}

func addCarrot(c *gin.Context) {
	owner := c.DefaultQuery("owner", "-1") // shortcut for c.Request.URL.Query().Get("lastname")
	rabbitOwner := extractPlayerFromStringID(owner)

	updatePlayerEnergy(&rabbitOwner)

	if rabbitOwner.Energy < carrotEnergyCost {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not enough energy to spawn carrot",
		})
	}

	xCoord, yCoord, err := extractCoordinates(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to extract carrot coordinates",
		})
	}

	rabbitOwner.Energy -= carrotEnergyCost

	newCarrot := Carrot{xCoordinate: xCoord, yCoordinate: yCoord, LifetimeEnd: time.Now().Add(carrotLifetime), EnergyCost: carrotEnergyCost, Creator: rabbitOwner}
	DB.Create(&newCarrot)
	DB.Save(&rabbitOwner)

	c.JSON(http.StatusOK, newCarrot)
}

func addBunny(c *gin.Context) {
	bunnyName := c.Param("bunnyName")
	owner := c.DefaultQuery("owner", "0") // shortcut for c.Request.URL.Query().Get("lastname")
	ownerId, _ := strconv.ParseUint(owner, 10, 64)

	newRabbit := Rabbit{Name: bunnyName, Owner: ownerId}
	result := DB.Create(&newRabbit)

	playerOwner := extractPlayerFromStringID(owner)

	spawnEvent := RabbitSpawn{CreatedRabbit: newRabbit, EnergyCost: 50, Creator: playerOwner}

	DB.Create(spawnEvent)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, newRabbit)
	} else {
		c.JSON(http.StatusBadRequest, newRabbit)
	}
}

//TODO error handling
func extractPlayerFromStringID(ownerIdStr string) RabbitOwner {
	ownerId, _ := strconv.ParseUint(ownerIdStr, 10, 64)

	var rabbitOwner RabbitOwner
	DB.Where("ID = ", ownerId).First(&rabbitOwner)

	return rabbitOwner
}

func updatePlayerEnergy(rabbitOwner *RabbitOwner) {
	timePassed := time.Now().Sub(rabbitOwner.LastEnergyCheck)
	rabbitOwner.Energy += timePassed.Seconds() * energyRechargeRate

	currentTime := time.Now()
	rabbitOwner.LastEnergyCheck = currentTime
}

func updatePlayerTimeChecked(rabbitOwner *RabbitOwner) {
}

func getBunniesSinceLastTime(c *gin.Context) {
	ownerParam := c.Param("ownerID")

	rabbitOwner := extractPlayerFromStringID(ownerParam)

	var rabbits []Rabbit
	DB.Where("created_at > ", rabbitOwner.LastCheck).Find(&rabbits)

	updatePlayerEnergy(&rabbitOwner)
	updatePlayerTimeChecked(&rabbitOwner)

	DB.Save(&rabbitOwner)

	c.JSON(http.StatusOK, rabbits)
}

func getBunnies(c *gin.Context) {
	var rabbits []Rabbit
	DB.Find(&rabbits)

	c.JSON(http.StatusOK, rabbits)
}

//func addThingRequest(c *gin.Context) {
//	code := c.Param("code")
//	price, _ := strconv.ParseUint(c.Query("price"), 10, 0)
//	amount, _ := strconv.ParseUint(c.DefaultQuery("amount", "0"), 10, 0)
//
//	newProduct := Product{
//		Code:   code,
//		Price:  price,
//		Amount: amount,
//	}
//
//	result := DB.Create(&newProduct)
//
//	if result.RowsAffected > 0 {
//		c.String(http.StatusOK, "Rows affected %i", result.RowsAffected)
//	} else {
//		c.String(http.StatusBadRequest, "Fail")
//	}
//}
//
//func getThings(c *gin.Context) {
//	var products []Product
//	DB.Find(&products)
//
//	c.JSON(http.StatusOK, products)
//}
