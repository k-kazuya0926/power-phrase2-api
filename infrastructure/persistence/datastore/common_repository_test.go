package datastore

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

func setup() {
	if err := godotenv.Load("../../../test.env"); err != nil {
		log.Fatal("Error loading test.env file")
	}
}

func teardown(db *gorm.DB) {
	db.DropTable(&model.Comment{})
	db.DropTable(&model.Post{})
	db.DropTable(&model.User{})
}
