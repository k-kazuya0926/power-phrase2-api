package datastore

import (
	"github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

func setup() {
}

func teardown(db *gorm.DB) {
	db.DropTable(&model.Favorite{})
	db.DropTable(&model.Comment{})
	db.DropTable(&model.Post{})
	db.DropTable(&model.User{})
}
