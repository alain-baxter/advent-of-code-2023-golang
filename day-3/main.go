package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"math"
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

	lineLength := 0
	indexOffset := 0
	parts := make(map[string][]int)
	symbols := make(map[int]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if lineLength == 0 {
			lineLength = len(text)
			log.Printf("Line length: %d", lineLength)
		}
		log.Println(text)
		parseSchematic(text, indexOffset, lineLength)
		lineParts, lineSymbols := parseSchematic(text, indexOffset, lineLength)

		maps.Copy(parts, lineParts)
		maps.Copy(symbols, lineSymbols)

		indexOffset += lineLength
	}

	log.Printf("Symbol Locations: %v", symbols)
	sumParts := 0
	gearTracker := make(map[int][]int)
	for key, locations := range parts {
		log.Printf("part: %s adjacent locations: %d", key, locations)
		for _, location := range locations {
			symbol, ok := symbols[location]
			if ok {
				part, _ := strconv.Atoi(strings.Split(key, ":")[1])
				sumParts += part

				if symbol == "*" {
					gearTracker[location] = append(gearTracker[location], part)
				}

				log.Println("Part Matched")
				break
			}
		}
	}

	log.Printf("Gear Connections: %v", gearTracker)
  sumGears := 0 
	for _, connections := range gearTracker {
		if len(connections) == 2 {
			sumGears += connections[0] * connections[1]
		}
	}

	log.Printf("Sum of parts: %d", sumParts)
	log.Printf("Sum of gear connections: %d", sumGears)
}

func parseSchematic(line string, indexOffset, lineLength int) (map[string][]int, map[int]string) {
	schematic := strings.Split(line, "")

	parts := make(map[string][]int)
	var partBuffer []int

	symbols := make(map[int]string)

	for lineIndex, item := range schematic {
		globalIndex := indexOffset + lineIndex
		isDigit := false
		
		if partDigit, err := strconv.Atoi(item); err == nil {
			partBuffer = append(partBuffer, partDigit)
			isDigit = true
		}

		partComplete, locModifier := partComplete(partBuffer, isDigit, lineIndex, lineLength)
		if partComplete {
			partNum := 0
			partLength := len(partBuffer)
			for j, num := range partBuffer {
				partNum += num * int(math.Pow10(partLength - j - 1))
			}
			partBuffer = nil

			var possibleLocations []int
			firstIndex := globalIndex + locModifier - partLength - 1
			lastIndex := globalIndex + locModifier

			// Add index before the number unless the number is at start of schematic line
			if lineIndex - partLength != 0 {
				possibleLocations = append(possibleLocations, firstIndex)
			} else {
				firstIndex += 1
			}

			// Add index after the number unless the number is at end of schematic line
			if lineIndex < lineLength - 1 {
				possibleLocations = append(possibleLocations, lastIndex)
			} else {
				lastIndex -= 1
			}
			
			// Add the indexes on the previous and next line
			for k := firstIndex; k <= lastIndex; k++ {
				possibleLocations = append(possibleLocations, k - lineLength)
				possibleLocations = append(possibleLocations, k + lineLength)
			}

			parts[fmt.Sprint(globalIndex) + ":" + fmt.Sprint(partNum)] = possibleLocations
		}

		if item != "." && !isDigit {
			symbols[globalIndex] = item
		}
	}

	return parts, symbols
}

func partComplete(partBuffer []int, isDigit bool, index, lineLength int) (bool, int) {
	// last entry in the line is a digit
	if isDigit && (index == lineLength - 1) {
		return true, 1
	}
	// first non-digit after a part
	if !isDigit && len(partBuffer) > 0 {
		return true, 0
	}
	return false, 0
}