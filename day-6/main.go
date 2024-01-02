package main

import (
	"alain-baxter/aoc-2023/day-6/race"
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

	var times, records []int
	var timebuffer, recordbuffer string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		fields := strings.Fields(text)
		ptr := &timebuffer
		if fields[0] == "Distance:" {
			ptr = &recordbuffer
		}
		for _, field := range fields[1:] {
			*ptr = *ptr + field
		}
	}

	time, _ := strconv.Atoi(timebuffer)
	times = append(times, time)

	record, _ := strconv.Atoi(recordbuffer)
	records = append(records, record)

	log.Printf("Times: %d Distance Records: %d", times, records)

	var winningOptions []int
	var waysToWin int = 1
	for i := 0; i < len(times); i++ {
		val := race.BeatRecord(times[i], records[i])
		winningOptions = append(winningOptions, val)
		waysToWin = waysToWin * val
	}
	log.Printf("Winning Options: %d", winningOptions)
	log.Printf("Number of ways to win: %d", waysToWin)
}
