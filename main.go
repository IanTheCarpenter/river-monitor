package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("UNABLE TO LOAD ENV FILE")
		return
	}
	fmt.Println(os.Getenv("TEST_VALUE"))

	fmt.Println("Hello, world")
	// client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))

}
