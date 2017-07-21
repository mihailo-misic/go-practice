package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Create your file with desired read/write permissions
	f, err := os.OpenFile("./go.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Set output of logs to f file
	log.SetOutput(f)

	log.Println("Heyooo")

	fmt.Println("DONE")
}
