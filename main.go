package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var db, err = gorm.Open(mysql.New(mysql.Config{
	DSN: "root:rootpass@tcp(127.0.0.1:3306)/gift-finder?charset=utf8&parseTime=True&loc=Local",
}), &gorm.Config{})

type GormModel struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-";sql:"index"`
}

type Gift struct {
	GormModel
	Name     string `gorm:"not null;size:255"`
	Gender   string `gorm:"not null;size:255"` // male|female|unisex
	AgeStart uint8  `gorm:"not null"`
	AgeEnd   uint8  `gorm:"not null"`
	Product  Product
	Images   []Image
	Hobbies  []Hobby `gorm:"many2many:gift_has_hobbies;"`
}

type Shop struct {
	GormModel
	Name string `gorm:"not null;size:255"`
}

type Product struct {
	GormModel
	GiftID uint    `gorm:"not null"`
	ShopID uint    `gorm:"not null"`
	Name   string  `gorm:"not null;size:255"`
	Price  float32 `gorm:"not null"`
	Url    string  `gorm:"not null"`
	Shop   Shop
}

type Image struct {
	GormModel
	GiftID uint   `gorm:"not null"`
	Url    string `gorm:"not null"`
}

type Hobby struct {
	GormModel
	Name string `gorm:"not null;size:255"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "{\"error\":\"Not Found\"}")
}

func hobbiesHandler(w http.ResponseWriter, r *http.Request) {
	var hobbies []Hobby
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "{\"error\":\"Only GET methods are supported\"}")
		return
	}

	result := db.Find(&hobbies)
	if result.Error != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": result.Error.Error()})
		return
	}

	json.NewEncoder(w).Encode(hobbies)
}

func main() {
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&Gift{}, &Shop{}, &Product{}, &Image{}, &Hobby{})

	http.HandleFunc("/", handler) // each request calls handler
	http.HandleFunc("/hobbies", hobbiesHandler)
	log.Fatal(http.ListenAndServe("localhost:8200", nil))
}
