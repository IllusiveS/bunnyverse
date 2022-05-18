package main

import (
	"gorm.io/gorm"
	"time"
)

type RabbitOwner struct {
	gorm.Model
	Name             string
	LastEnergyCheck  time.Time
	LastBunniesCheck time.Time
	LastCarrotsCheck time.Time
	Energy           float64
	Rabbits          []Rabbit `gorm:"foreignKey:ID;references:ID"`
}

type Carrot struct {
	gorm.Model
	LifetimeEnd time.Time
	xCoordinate float64
	yCoordinate float64
	EnergyCost  float64
	Creator     RabbitOwner
}

type RabbitSpawn struct {
	gorm.Model
	CreatedRabbit Rabbit
	EnergyCost    float64
	Creator       RabbitOwner
}

type Rabbit struct {
	gorm.Model
	Name       string
	Fluffiness uint32 //0-1000
	EarLength  float64
	Dots       uint64
	Bounces    uint64
	R          uint32
	G          uint32
	B          uint32
	Owner      uint64
	Parents    []Rabbit `gorm:"foreignKey:ID;references:ID"`
}
