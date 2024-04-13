package main

import (
	"log"
	"os"
)

func writeDataFile(filename string, data []byte) {
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
