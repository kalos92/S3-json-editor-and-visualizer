package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(T *testing.T) {
	err := godotenv.Load("source.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	by := `"{\"ciao\":1}"`

	by2, err := strconv.Unquote(by)

	if err != nil {
		panic(err)
	}

	b := []byte(by2)

	var f interface{}
	json.Unmarshal(b, &f)

	myMap := f.(map[string]interface{})

	fmt.Println(myMap["key"])
}
