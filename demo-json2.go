package main

import (
	"encoding/json"
	"fmt"
)

type employee struct {
	ID    int
	Name  string
	Tel   string
	Email string
}

func main() {
	e := employee{}
	err := json.Unmarshal([]byte(`{"ID": 101, "Name": "buncha1", "Tel": "0215485534","Email":"rr@gmail.com"}`), &e)
	if err != nil {
		panic(err)
	}
	fmt.Println(e)

	_data, _ := json.Marshal(&e)
	fmt.Println(string(_data))

	fmt.Println(e.Email)

}
