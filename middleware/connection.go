package middleware

import (
	"database/sql"
	"fmt"

	"github.com/spf13/viper"
)

var (
	host     = viper.GetString("dbConfig.host")
	port     = viper.GetString("dbConfig.port")
	user     = viper.GetString("dbConfig.user")
	password = viper.GetString("dbConfig.password")
	dbname   = viper.GetString("dbConfig.dbname")
)

var DB *sql.DB

func CreateConnection() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	DB = db

	fmt.Println("Successfully connected!")
	return db
}
