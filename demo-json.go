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
	_data, _ := json.Marshal(&employee{101, "buncha rr", "012547788", "buncha.rrr@gmail.com"})
	fmt.Println(string(_data))

}
