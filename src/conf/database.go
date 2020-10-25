package conf

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/k-kazuya0926/power-phrase2-api/domain/model"
)

// NewDBConnection 新規データベースコネクションを取得します.
func NewDBConnection() *gorm.DB {
	return getMysqlConnection()
}

func getMysqlConnection() *gorm.DB {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		Current.Database.User,
		Current.Database.Password,
		Current.Database.Host,
		Current.Database.Port,
		Current.Database.Database,
	)

	connection, err := gorm.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	err = connection.DB().Ping()
	if err != nil {
		panic(err)
	}

	connection.LogMode(true)
	connection.DB().SetMaxIdleConns(10)
	connection.DB().SetMaxOpenConns(20)

	connection.Set("gorm:table_options", "ENGINE=InnoDB")

	connection.AutoMigrate(&model.User{})
	connection.AutoMigrate(&model.Post{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")

	return connection
}
