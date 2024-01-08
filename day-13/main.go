package main

import (
	"bufio"
	"log"
	"os"
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

	var rowBuffer [][]string
	summaryCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		// Reset and detect when a section is complete
		if len(strings.TrimSpace(text)) == 0 {
			summaryCount += detectMirrors(rowBuffer)
			rowBuffer = [][]string{}
			continue
		}
		log.Println(text)
		rowBuffer = append(rowBuffer, getSymbols(text))
	}

	// Get the last mirror
	if (len(rowBuffer) != 0) {
		summaryCount += detectMirrors(rowBuffer)
	}

	log.Printf("Mirror summarization: %d", summaryCount)
}

func getSymbols(text string) []string {
	return strings.Split(text, "")
}

func detectMirrors(rows [][]string) int {
	// Horizontal Mirrors
	count, found := checkForMirror(rows)
	if found {
		log.Printf("Horizontal mirror found, returning %d", count * 100)
		return count * 100
	} else {
		// Vertical Mirrors
		count, found := checkForMirror(switchToColumns(rows))
		if found {
			log.Printf("Vertical mirror found, returning %d", count)
			return count
		}
	}
	log.Fatal("No mirror detected")
	return 0
}

func checkForMirror(sections [][]string) (int, bool) {
	for index, symbols := range sections {
		matched := false
		smudged := false
		// Leave if there is no next
		if index == len(sections) - 1 {
			break
		}

		// This is not a mirror if it does not match the next entry
		next := sections[index + 1]
		matched, smudged = compareSections(symbols, next, smudged)
		if !matched {
			continue
		}

		log.Printf("Duplicate %s at %d with smudge: %t", symbols, index, smudged)

		// Check for mirror
		valid := true
		upwards := index + 1

		for j := 1; j >= 0; j++ {
			first := index - j
			second := upwards + j
			// If we go out of range end
			if first < 0 || second >= len(sections) {
				break
			}

			log.Printf("Comparing %d:%s with %d:%s", first, sections[first], second, sections[second])
			valid, smudged = compareSections(sections[first], sections[second], smudged)
			log.Printf("Valid: %t Smudge: %t", valid, smudged)

			// If we detect a difference end
			if !valid {
				break
			}
		}

		// Only use smudged answers
		if valid && smudged {
			return index + 1, true
		}
	}

	return -1, false
}

func compareSections(section1, section2 []string, smudged bool) (bool, bool) {
	matched := true

	smudgeCount := 0
	if smudged {
		smudgeCount = 1
	}

	for i := range section1 {
		first := string(section1[i])
		second := string(section2[i])
		if first != second {
			if smudgeCount == 0 {
				smudgeCount = 1
			} else {
				matched = false
			}
		}
	}

	return matched, smudgeCount == 1
}

func switchToColumns(rows [][]string) [][]string {
	cols := make([][]string, len(rows[0]))
	for _, symbols := range rows {
		for index, symbol := range symbols {
			cols[index] = append(cols[index], symbol)
		}
	}
	return cols
}