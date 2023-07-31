package database

import (
	"database/sql"
	"fmt"
	"log"
)

type postgresDatabase struct {
	db *sql.DB
}

func (x *postgresDatabase) BuildDatabaseUri(dbData DatabaseInfo) string {
	return "postgres://" + dbData.DbUsername + ":" + dbData.DbPass + "@" + dbData.DbHost + ":" + dbData.DbPort + "/" + dbData.DbName
}

func (x *postgresDatabase) ConnectDatabase(dbData DatabaseInfo) error {
	uri := x.BuildDatabaseUri(dbData)
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	x.db = db
	return nil
}

func (x *postgresDatabase) GetDatabaseConnection() *sql.DB {
	if x.db != nil {
		return x.db
	}
	fmt.Println("No database connection")
	return nil
}
