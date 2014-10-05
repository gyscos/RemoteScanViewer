package main

// Could have called a go-sane binding to directly query the printer...
// ... or I can just use a bash script.

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func findFirstFree(directory string, prefix string) int {
	iMax := 99999
	for i := 0; i < iMax; i++ {
		filename := directory + prefix + strconv.Itoa(i)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return i
		}
	}
	log.Println("Error: no more free number")
	return iMax
}

func scanImage(directory string) {
	namePrefix := time.Now().Format("2006-01-02 ")
	i := findFirstFree(directory, namePrefix)
	filename := namePrefix + strconv.Itoa(i)
	err := exec.Command("scripts/scan.sh", filename).Run()
	if err != nil {
		log.Println("Error:", err)
	}
}
