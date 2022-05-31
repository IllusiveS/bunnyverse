package test_rest

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=db user=postgres password=postgres dbname=bunnycalypse port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Rabbit{})
	database.AutoMigrate(&RabbitOwner{})
	database.AutoMigrate(&Carrot{})

	DB = database
}
