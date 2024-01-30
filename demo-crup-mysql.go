package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Course struct {
	id       int     `่json: "id"`
	name     string  `่json: "name"`
	price    float64 `่json: "price"`
	imageURL string  `่json: "image_url"`
}

func createTable(db *sql.DB) {
	query := `CREATE TABLE users (
		id INT AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created DATETIME,
		PRIMARY KEY (id)
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func insert(db *sql.DB) {
	var username string
	var password string
	fmt.Scan(&username)
	fmt.Scan(&password)
	created := time.Now()
	result, err := db.Exec(`insert into users (username,password,created) value (?,?,?)`, username, password, created)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	fmt.Println("id : ", id)
}

func delete(db *sql.DB) {
	var deleteID int
	fmt.Scan(&deleteID)
	_, err := db.Exec(`delete from users where id = ?`, deleteID)
	if err != nil {
		log.Fatal(err)
	}
}

func query2(db *sql.DB) ([]Course, error) {
	// var (
	// 	id         int
	// 	coursname  string
	// 	price      float64
	// 	instructor string
	// )
	// for {
	// 	var inputID int
	// 	fmt.Scan(&inputID)
	// 	query := "select id,coursname,price,instructor from onlinecourse where id = ?"
	// 	if err := db.QueryRow(query, inputID).Scan(&id, &coursname, &price, &instructor); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(id, coursname, price, instructor)
	// }
	results, err := db.Query(`SELECT * FROM courseonline`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	courses := make([]Course, 0)
	for results.Next() {
		var course Course
		results.Scan(&course.id, &course.name, &course.price, &course.imageURL)
		courses = append(courses, course)
	}
	return courses, nil
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/dbdev")
	if err != nil {
		fmt.Println("Failed to connect")
	} else {
		fmt.Println("connect successfully")
	}
	list, err := query2(db)
	fmt.Println("test:", list)
	//createTable(db)
	//insert(db)
	//delete(db)
}
