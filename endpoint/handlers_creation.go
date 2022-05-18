package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func addBunny(c *gin.Context) {
	//TODO move param to query
	bunnyName := c.Param("bunnyName")
	owner := c.DefaultQuery("owner", "0") // shortcut for c.Request.URL.Query().Get("lastname")
	ownerId, _ := strconv.ParseUint(owner, 10, 64)

	newRabbit := Rabbit{Name: bunnyName, Owner: ownerId}
	result := DB.Create(&newRabbit)

	playerOwner, _ := extractPlayerFromStringID(owner)

	spawnEvent := RabbitSpawn{CreatedRabbit: newRabbit, EnergyCost: 50, Creator: *playerOwner}

	DB.Create(spawnEvent)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, newRabbit)
	} else {
		c.JSON(http.StatusBadRequest, newRabbit)
	}
}

func addCarrot(c *gin.Context) {
	owner := c.DefaultQuery("owner", "-1")
	rabbitOwner, _ := extractPlayerFromStringID(owner)

	updatePlayerEnergy(rabbitOwner)

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

	newCarrot := Carrot{xCoordinate: xCoord, yCoordinate: yCoord, LifetimeEnd: time.Now().Add(carrotLifetime), EnergyCost: carrotEnergyCost, Creator: *rabbitOwner}
	DB.Create(&newCarrot)
	DB.Save(&rabbitOwner)

	c.JSON(http.StatusOK, newCarrot)
}

func createUser(name string) *RabbitOwner {
	now := time.Now()
	owner := RabbitOwner{Name: name, LastCarrotsCheck: now, LastBunniesCheck: now, LastEnergyCheck: now, Energy: 100}
	DB.Create(&owner)

	return &owner
}

func connectUser(c *gin.Context) {
	player, err := extractPlayerFromNameContext(c)

	if err != nil {
		if err.Error() != "failed to find player" {
			c.JSON(http.StatusBadRequest, err)
		} else {
			owner := c.DefaultQuery("owner", "")
			player = createUser(owner)
		}
	}

	c.JSON(http.StatusOK, player)
}
