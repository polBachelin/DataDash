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

func (x *postgresDatabase) ConnectDatabase(dbData DatabaseInfo) *sql.DB {
	uri := x.BuildDatabaseUri(dbData)
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	x.db = db
	return db
}

func (x *postgresDatabase) GetDatabaseConnection() *sql.DB {
	if x.db != nil {
		return x.db
	}
	fmt.Println("No database connection")
	return nil
}
