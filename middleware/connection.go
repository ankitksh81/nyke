package middleware

import (
	"database/sql"

	"github.com/ankitksh81/nyke/logger"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sql.DB

func CreateConnection() *sql.DB {
	logger.Log.Info("Connecting to database...")
	var ConStr = viper.GetString("dbConfig.constr")
	db, err := sql.Open("postgres", ConStr)
	if err != nil {
		logger.Log.Info("Could not create connection to the database" + err.Error())
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Log.Info("Error connecting to the database" + err.Error())
		panic(err)
	}

	DB = db

	logger.Log.Info("Successfully connected to the database!")
	return db
}
