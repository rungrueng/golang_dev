package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func query(db *sql.DB) {
	var (
		id         int
		coursname  string
		price      float64
		instructor string
	)
	for {
		var inputID int
		fmt.Scan(&inputID)
		query := "select id,coursname,price,instructor from onlinecourse where id = ?"
		if err := db.QueryRow(query, inputID).Scan(&id, &coursname, &price, &instructor); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, coursname, price, instructor)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/dbdev")
	if err != nil {
		fmt.Println("Failed to connect")
	} else {
		fmt.Println("connect successfully")
	}
	query(db)
}
