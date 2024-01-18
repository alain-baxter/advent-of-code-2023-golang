package utils

import (
	"bufio"
	"log"
	"os"
)

func LoadFile(args []string) []string {
	var filepath string

	if len(args) > 1 {
		filepath = args[1]
	} else {
		log.Fatal("Need to pass file path as argument")
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		lines = append(lines, text)
	}

	return lines
}