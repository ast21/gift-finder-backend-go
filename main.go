package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db, err = gorm.Open(mysql.New(mysql.Config{
	DSN: "root:rootpass@tcp(127.0.0.1:3306)/gift-finder?charset=utf8&parseTime=True&loc=Local",
}), &gorm.Config{})

type Gift struct {
	gorm.Model
	Name     string `gorm:"not null;size:255"`
	Gender   string `gorm:"not null;size:255"` // male|female|unisex
	AgeStart uint8  `gorm:"not null"`
	AgeEnd   uint8  `gorm:"not null"`
	Product  Product
	Images   []Image
	Hobbies  []Hobby `gorm:"many2many:gift_has_hobbies;"`
}

type Shop struct {
	gorm.Model
	Name string `gorm:"not null;size:255"`
}

type Product struct {
	gorm.Model
	GiftID uint    `gorm:"not null"`
	ShopID uint    `gorm:"not null"`
	Name   string  `gorm:"not null;size:255"`
	Price  float32 `gorm:"not null"`
	Url    string  `gorm:"not null"`
	Shop   Shop
}

type Image struct {
	gorm.Model
	GiftID uint   `gorm:"not null"`
	Url    string `gorm:"not null"`
}

type Hobby struct {
	gorm.Model
	Name string `gorm:"not null;size:255"`
}

func main() {
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&Gift{}, &Shop{}, &Product{}, &Image{}, &Hobby{})
	fmt.Printf("Gift finder")
}
