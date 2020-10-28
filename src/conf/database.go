package conf

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

// NewDBConnection 新規データベースコネクションを取得します.
func NewDBConnection() *gorm.DB {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		Current.Database.User,
		Current.Database.Password,
		Current.Database.Host,
		Current.Database.Port,
		Current.Database.Database,
	)
	return getMysqlConnection(connectionString)
}

func NewTestDBConnection() *gorm.DB {
	// TODO 定義移動
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		"kazuya",
		"kazuya",
		"localhost",
		"3306",
		"power-phrase2-test",
	)
	return getMysqlConnection(connectionString)
}

func getMysqlConnection(connectionString string) *gorm.DB {
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(20)

	db.Set("gorm:table_options", "ENGINE=InnoDB")

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	return db
}
