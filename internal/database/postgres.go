package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type postgresDatabase struct {
	db *sql.DB
}

func (x *postgresDatabase) BuildDatabaseUri(dbData DatabaseInfo) string {
	return "postgres://" + dbData.DbUsername + ":" + dbData.DbPass + "@" + dbData.DbHost + ":" + dbData.DbPort + "/" + dbData.DbName + "?sslmode=disable"
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

func (x *postgresDatabase) ExecuteQuery(query interface{}) (interface{}, error) {
	log.Printf("%s", fmt.Sprintf("%s", query))
	res, err := x.db.Query(fmt.Sprintf("%s", query))
	if err != nil {
		fmt.Println("Error executing query: ", err)
		return nil, err
	}
	return res, nil
}
