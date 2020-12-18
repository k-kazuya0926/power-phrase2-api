package conf

import (
	"log"

	"github.com/joho/godotenv"
)

// NewConfig プロジェクトのコンフィグ設定をロードする。
func NewConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	return
}
