package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var db, err = gorm.Open(postgres.New(postgres.Config{
	DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env("DB_HOST"),
		env("DB_PORT"),
		env("DB_USER"),
		env("DB_PASS"),
		env("DB_NAME"),
	),
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
	Hobbies  []*Hobby `gorm:"many2many:gift_has_hobbies;"`
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
	Name  string  `gorm:"not null;size:255"`
	Gifts []*Gift `gorm:"many2many:gift_has_hobbies;"`
}

type GiftRequestBody struct {
	HobbyIds []int
	Gender   string
	Age      string
}

func env(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
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
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": result.Error.Error()})
		return
	}

	json.NewEncoder(w).Encode(hobbies)
}

func giftsHandler(w http.ResponseWriter, r *http.Request) {
	var hobbies []Hobby
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "{\"error\":\"Only POST methods are supported\"}")
		return
	}

	HobbyIds, err := stringToNumbers(r.URL.Query().Get("hobby_ids"))
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	req := GiftRequestBody{
		HobbyIds,
		r.URL.Query().Get("gender"),
		r.URL.Query().Get("age"),
	}

	if len(req.HobbyIds) > 0 {
		result := db.Preload("Gifts", "gender = ? AND age_start <= ? AND age_end >= ?", req.Gender, req.Age, req.Age).Preload("Gifts.Product").Preload("Gifts.Product.Shop").Preload("Gifts.Images").Preload("Gifts.Hobbies").Find(&hobbies, "id", req.HobbyIds)
		if result.Error != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": result.Error.Error()})
			return
		}
	} else {
		result := db.Preload("Gifts", "gender = ? AND age_start <= ? AND age_end >= ?", req.Gender, req.Age, req.Age).Preload("Gifts.Product").Preload("Gifts.Product.Shop").Preload("Gifts.Images").Preload("Gifts.Hobbies").Find(&hobbies)
		if result.Error != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(map[string]string{"error": result.Error.Error()})
			return
		}
	}

	gifts := funk.FlatMap(hobbies, func(x Hobby) []*Gift {
		return x.Gifts
	})
	gifts = funk.Uniq(gifts)

	json.NewEncoder(w).Encode(gifts)
}

func stringToNumbers(str string) ([]int, error) {
	if str == "" {
		return []int{}, err
	}
	split := strings.Split(str, ",")
	var numbers []int
	for i := 0; i < len(split); i++ {
		number, err := strconv.Atoi(split[i])
		if err != nil {
			return []int{}, err
		}
		numbers = append(numbers, number)
	}
	return numbers, err
}

func main() {
	if err != nil {
		panic("Failed to connect database")
	}

	db.AutoMigrate(&Gift{}, &Shop{}, &Product{}, &Image{}, &Hobby{})

	http.HandleFunc("/", handler) // each request calls handler
	http.HandleFunc("/hobbies", hobbiesHandler)
	http.HandleFunc("/gifts", giftsHandler)
	log.Fatal(http.ListenAndServe("localhost:8200", nil))
}
