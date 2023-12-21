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
	sumPossible := 0
	sumPower := 0
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		game, red, green, blue := parseGame(text)
		impossible := false
		if red > 12 {
			impossible = true
		}
		if green > 13 {
			impossible = true
		}
		if blue > 14 {
			impossible = true
		}

		if !impossible {
			sumPossible += game
		}
		sumPower += red * green * blue
	}

	log.Printf("Sum of possible games: %d", sumPossible)
	log.Printf("Sum of dice power for all games: %d", sumPower)
}

func parseGame(s string) (int, int, int, int) {
	rawGame := strings.Split(s, ":")
	gameNum, _ := strconv.Atoi(strings.Split(rawGame[0], " ")[1])
	maxRed := 0
	maxGreen := 0
	maxBlue := 0

	for _, set := range strings.Split(rawGame[1], ";") {
		for _, dice := range strings.Split(set, ",") {
			switch {
			case strings.Contains(dice, "red"):
				maxRed = replaceMax(maxRed, dice)
			case strings.Contains(dice, "green"):
				maxGreen = replaceMax(maxGreen, dice)
			case strings.Contains(dice, "blue"):
				maxBlue = replaceMax(maxBlue, dice)
			}
		}
	}
	
	return gameNum, maxRed, maxGreen, maxBlue
}

func replaceMax(max int, dice string) int {
	num, _ := strconv.Atoi(strings.Split(strings.TrimSpace(dice), " ")[0])
	if num > max {
		return num
	} else {
		return max
	}
}