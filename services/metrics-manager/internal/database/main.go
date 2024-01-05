package database

import (
	"database/sql"
	"log"
	"metricsTerrarium/lib"
)

type Db struct {
	Connection *sql.DB
}

type DbProperties struct {
	Config *lib.Config
}

func (db Db) Start(properties DbProperties) Db {
	connStr := "user=" + properties.Config.MetricsDataBaseUser + " password=" + properties.Config.MetricsDataBasePassword + " dbname=" + properties.Config.MetricsDataBase + " sslmode=disable"
	dbConnection, err := sql.Open("postgres", connStr)
	db.Connection = dbConnection

	if err != nil {
		log.Fatalf("Error during db connection creation. Err: %s", err)
	} else {
		log.Printf("Connection to database created succesfully")
	}

	defer db.Connection.Close()

	return db
}
