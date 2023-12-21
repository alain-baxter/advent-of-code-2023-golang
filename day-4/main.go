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
	sumScore := 0
	gameNum := 1
	gameTracker := make(map[int]int)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		score, matches := scoreGame(text)
		log.Printf("Score: %d", score)
		log.Printf("Matches: %d", matches)
		sumScore += score

		timesPlayed := 1
		v, ok := gameTracker[gameNum]
		if ok {
			gameTracker[gameNum] = v + 1
			timesPlayed = gameTracker[gameNum]
		} else {
			gameTracker[gameNum] = 1
		}
		log.Printf("Card %d Played %d times", gameNum, timesPlayed)

		for i := 1; i <= matches; i++ {
			v, ok := gameTracker[gameNum + i]
			if ok {
				gameTracker[gameNum + i] = v + timesPlayed
			} else {
				gameTracker[gameNum + i] = timesPlayed
			}
		}
		gameNum++
	}

	sumGames := 0
	for _, count := range gameTracker {
		sumGames += count
	}

	log.Printf("Sum of scores: %d", sumScore)
	log.Printf("Number of games: %d", sumGames)
}

func scoreGame(game string) (int, int) {
	score := 0
	matches := 0
	first := true
	gameData := strings.Split(game, ":")
	numberData := strings.Split(gameData[1], "|")

	winningNumbers := make(map[int]bool)
	var numbers []int
	for _, num := range strings.Split(numberData[0], " ") {
		if (len(strings.TrimSpace(num))) == 0 {
			continue
		}
		value, _ := strconv.Atoi(strings.TrimSpace(num))
		winningNumbers[value] = true
	}
	for _, num := range strings.Split(numberData[1], " ") {
		if (len(strings.TrimSpace(num))) == 0 {
			continue
		}
		value, _ := strconv.Atoi(strings.TrimSpace(num))
		numbers = append(numbers, value)
	}

	for _, num := range numbers {
		_, ok := winningNumbers[num]
		if ok && first {
			score = 1
			first = false
			matches++
		} else if ok {
			score = score * 2
			matches++
		}
	}

	return score, matches
}
