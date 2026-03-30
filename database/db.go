package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	fmt.Println("Connect() called") // debug line
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/nooboj")
	if err != nil {
		panic("SQL Open failed: " + err.Error())
	}

	err = DB.Ping()
	if err != nil {
		panic("Ping failed: " + err.Error())
	}

	fmt.Println("Database connected")
}
