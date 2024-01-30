package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Course1 struct {
	Id       int     `่json: "id"`
	Name     string  `่json: "name"`
	Price    float64 `่json: "price"`
	ImageURL string  `่json: "imge_url"`
}

var Db *sql.DB

const coursePath = "course"
const basePath = "/api"

func getCourseList() ([]Course1, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := Db.QueryContext(ctx, `SELECT id,name,price,imge_url FROM courseonline`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	courses := make([]Course1, 0)
	for results.Next() {
		var course Course1
		results.Scan(&course.Id, &course.Name, &course.Price, &course.ImageURL)
		courses = append(courses, course)
	}
	return courses, nil
}

func insertCourse(course Course1) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.ExecContext(ctx, `insert into courseonline (id,name,price,imge_url) value (?,?,?,?)`, course.Id, course.Name, course.Price, course.ImageURL)
	if err != nil {
		log.Panicln(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Panicln(err.Error())
		return 0, err
	}

	return int(insertID), nil
}

func getCourseById(id int) (*Course1, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "select id,name,price,imge_url from courseonline where id = ?"
	row := Db.QueryRowContext(ctx, query, id)

	course := &Course1{}
	err := row.Scan(&course.Id, &course.Name, &course.Price, &course.ImageURL)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Panicln(err)
		return nil, err
	}
	return course, nil

}

func removeCourse(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "delete from courseonline where id = ?"
	_, err := Db.ExecContext(ctx, query, id)

	if err != nil {
		log.Panicln(err.Error())
		return err
	}
	return nil

}

func handlerCourse(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(&courseList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var newCourse Course1
		err := json.NewDecoder(r.Body).Decode(&newCourse)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		Id, err := insertCourse(newCourse)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"course Id ":%d}`, Id)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlerCourseByID(w http.ResponseWriter, r *http.Request) {

	urlPathSegment := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	if len(urlPathSegment[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1]) // string to int
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		course, err := getCourseById(ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodDelete:
		err := removeCourse(ID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func coursMiddleware(Handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Authorization,X-Custom-Header, Upgrade-Insecure-Requests,x-requested-with")
		Handler.ServeHTTP(w, r)

	})
}

func SetRoutes(apiPath string) {
	courseHandler := http.HandlerFunc(handlerCourse)
	courseHandlerByID := http.HandlerFunc(handlerCourseByID)
	http.Handle(fmt.Sprintf("%s/%s", apiPath, coursePath), coursMiddleware(courseHandler))
	http.Handle(fmt.Sprintf("%s/%s/", apiPath, coursePath), coursMiddleware(courseHandlerByID))

}

func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/dbdev")
	if err != nil {
		fmt.Println("Failed to connect")
	} else {
		fmt.Println("connect successfully")
	}
	fmt.Println(Db)
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func main() {
	SetupDB()
	SetRoutes(basePath)
	http.ListenAndServe(":5000", nil)
	//log.Fatal(http.ListenAndServe(":80", nil))
}
