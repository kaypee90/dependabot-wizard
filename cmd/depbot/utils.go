package main

import (
	"fmt"
	"log"
	"os"
)

func displayIntroductoryText() {
	green := "\033[32m"
	reset := "\033[0m"
	fmt.Printf("%sDepbot wizard will help you configure dependabot in your project%s\n", green, reset)
}

func writeDataToFile(filename string, data []byte) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error creating file %s: %v", filename, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("error writing to file %s: %v", filename, err)
	}

	log.Printf("Data has been written to %s", filename)
}
