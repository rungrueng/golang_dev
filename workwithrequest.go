package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Course struct {
	ID         int     `json: "id"`
	Name       string  `json: "name"`
	Price      float64 `json: "price"`
	Instructor string  `json: "instructor"`
}

var CoursList []Course

func init() {
	// CourseJSON := `[
	// 	{"id": 1,"name":"python","price":2500,"instructor":"cha dev"},
	// 	{"id": 2,"name":"JavaScript","price":3500,"instructor":"cha dev"},
	// 	{"id": 3,"name":"React","price":4500,"instructor":"cha dev"}
	// 	]`
	CourseJSON := `[]`

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

func courseHandler(w http.ResponseWriter, r *http.Request) {

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

func main() {
	http.HandleFunc("/course", courseHandler)
	http.ListenAndServe(":80", nil)

}
