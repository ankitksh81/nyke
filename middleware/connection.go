package middleware

import (
	"database/sql"
	"time"

	"github.com/ankitksh81/nyke/logger"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// global variable for database instance
var DB *sql.DB

func CreateConnection() *sql.DB {
	logger.Log.Info("Connecting to database...")
	var ConStr = viper.GetString("dbConfig.constr") // Don't write this line outside of this function(viper initialization issues)

	db, err := sql.Open("postgres", ConStr)
	if err != nil {
		logger.Log.Error("Could not create connection to the database" + err.Error())
		panic(err)
	}
	// defer db.Close()

	/* Database connections config */
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		logger.Log.Error("Error connecting to the database" + err.Error())
		panic(err)
	}

	DB = db

	logger.Log.Info("Successfully connected to the database!")
	return db
}
