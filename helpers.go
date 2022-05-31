package test_rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

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

func extractPlayerFromNameContext(c *gin.Context) (*RabbitOwner, error) {
	owner := c.DefaultQuery("owner", "")
	if owner == "" {
		return nil, errors.New("failed to extract owner name from query")
	}
	return extractPlayerFromStringName(owner)
}

func extractPlayerFromStringName(ownerName string) (*RabbitOwner, error) {
	if ownerName == "" {
		return nil, errors.New("empty user name")
	}

	var rabbitOwner RabbitOwner

	result := DB.Where("Name = ?", ownerName).First(&rabbitOwner)

	if result.Error == nil {
		return &rabbitOwner, nil
	} else {
		return nil, errors.New("failed to find player")
	}
}

func extractPlayerFromIDContext(c *gin.Context) (*RabbitOwner, error) {
	owner := c.DefaultQuery("owner", "-1")
	if owner == "-1" {
		return nil, errors.New("failed to extract owner ID from query")
	}
	return extractPlayerFromStringID(owner)
}

func extractPlayerFromStringID(ownerIdStr string) (*RabbitOwner, error) {
	ownerId, err := strconv.ParseUint(ownerIdStr, 10, 64)

	if err != nil {
		return nil, err
	}

	var rabbitOwner RabbitOwner
	DB.Where("ID = ?", ownerId).First(&rabbitOwner)

	return &rabbitOwner, nil
}

func updatePlayerEnergy(rabbitOwner *RabbitOwner) {
	timePassed := time.Now().Sub(rabbitOwner.LastEnergyCheck)
	rabbitOwner.Energy += timePassed.Seconds() * energyRechargeRate

	currentTime := time.Now()
	rabbitOwner.LastEnergyCheck = currentTime
}

func updatePlayerCarrotTimeChecked(rabbitOwner *RabbitOwner) {
	currentTime := time.Now()
	rabbitOwner.LastCarrotsCheck = currentTime
}
