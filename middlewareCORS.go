package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Course struct {
	ID         int     `json: "id"`
	Name       string  `json: "name"`
	Price      float64 `json: "price"`
	Instructor string  `json: "instructor"`
}

var CoursList []Course

func init() {
	CourseJSON := `[
		{"id": 1,"name":"python","price":2500,"instructor":"cha dev"},
		{"id": 2,"name":"JavaScript","price":3500,"instructor":"cha dev"},
		{"id": 3,"name":"React","price":4500,"instructor":"cha dev"}
		]`
	//CourseJSON := `[]`

	err := json.Unmarshal([]byte(CourseJSON), &CoursList)
	if err != nil {
		log.Fatal(err)
	}
}

func getNextID() int {
	highestID := -1
	if len(CoursList) != 0 {
		for _, course := range CoursList {
			if highestID < course.ID {
				highestID = course.ID
			}
		}
	} else {
		highestID = 0
	}

	return highestID + 1
}

func findID(ID int) (*Course, int) {
	for i, course := range CoursList {
		if course.ID == ID {
			return &course, i
		}
	}
	return nil, 0

}

func courseByID(w http.ResponseWriter, r *http.Request) {

	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1]) // string to int
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no course with id %d", ID), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		courseJSON, err := json.Marshal(course)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Cotent-type", "application/json")
		w.Write(courseJSON)

	case http.MethodPut:
		var updateCourse Course
		BodyByte, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(BodyByte, &updateCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updateCourse.ID != ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		course = &updateCourse
		CoursList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func courseHandlerall(w http.ResponseWriter, r *http.Request) {

	courseJSON, err := json.Marshal(CoursList)
	switch r.Method {
	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)

	case http.MethodPost:
		var newCourse Course
		BodyByte, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(BodyByte, &newCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newCourse.ID = getNextID()
		CoursList = append(CoursList, newCourse)
		newCourseJSON, err := json.Marshal(newCourse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(newCourseJSON)
		return

	}

}

func corsMiddlewareHandler(Handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,PATCH")
		w.Header().Add("Access-Control-Allow-Headers", "Accept,Content-Type,Content-Length,Authorization,X-Custom-Header, Upgrade-Insecure-Requests,x-requested-with")
		Handler.ServeHTTP(w, r)

	})
}

func main() {

	courseByItem := http.HandlerFunc(courseByID)
	courseList := http.HandlerFunc(courseHandlerall)

	http.Handle("/course/", corsMiddlewareHandler(courseByItem))
	http.Handle("/course", corsMiddlewareHandler(courseList))
	http.ListenAndServe(":5000", nil)

}
