package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var filepath string

	args := os.Args
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
	sum := 0
	history := 0
	for scanner.Scan() {
		text := scanner.Text()
		
		var seq []int
		for _, val := range strings.Fields(text) {
			ival, _ := strconv.Atoi(val)
			seq = append(seq, ival)
		}
		log.Println(seq)

		next, hist := getPredictedValue(seq)
		sum += next
		history += hist
		log.Printf("Next prediction for %d is %d", seq, next)
		log.Printf("History prediction for %d is %d", seq, hist)
	}
	log.Printf("Sum of next predictions: %d", sum)
	log.Printf("Sum of history predictions: %d", history)
}

func getPredictedValue(seq []int) (int, int) {
	zeroes := true
	for _, val := range seq {
		if val != 0 {
			zeroes = false
		}
	}

	if zeroes {
		return 0, 0
	}

	var nextSeq []int
	for i := 1; i < len(seq); i++ {
		nextSeq = append(nextSeq, seq[i] - seq[i - 1])
	}
	log.Print(nextSeq)

	next, hist := getPredictedValue(nextSeq)
	return seq[len(seq) - 1] + next, seq[0] - hist
}