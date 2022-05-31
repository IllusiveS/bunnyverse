package test_rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var energyRechargeRate = 1.0
var carrotEnergyCost = 20.0
var carrotLifetime = time.Second * 30
var bunnyEnergyCost = 100.0

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
