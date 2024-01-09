package main

import (
	"bufio"
	"fmt"
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

	var platform [][]string
	rowCount := 0
	cycles := 1000000000
	logs := false
	scanner := bufio.NewScanner(file)
	log.Println("Platform:")
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		platform = append(platform, strings.Split(text, ""))
		rowCount++
	}

	//tiltPlatform(platform)
	load := cyclePlatformMultiple(platform, cycles, logs)
	log.Printf("Total Load: %d", load)
}

func tiltPlatform(platform [][]string) {
	// The Row
	for i, symbols := range platform {
		//log.Printf("Working on %d:%s", i, symbols)
		// The columns in the row
		for j, symbol := range symbols {
			//log.Printf("Tilting at %d:%s", j, symbol)
			if symbol == "O" || symbol == "#" {
				continue
			}

			// We have a ".", find the next "O" in the column
			for k := i + 1; k < len(platform); k++ {
				check := platform[k][j]
				//log.Printf("Rolling from %d:%s", k, check)
				// Can't move "#" or anything after it
				if check == "#" {
					break
				}

				// Move "O" to this row
				if check == "O" {
					platform[k][j] = "."
					platform[i][j] = "O"
					break
				}
			}
		}
	}
}

func getLoad(platform [][]string) int {
	rowCount := len(platform)
	load := 0
	for row, symbols := range platform {
		load += (rowCount - row) * countSymbols(symbols, "O")
	}
	return load
}

func cyclePlatformMultiple(platform [][]string, times int, logs bool) int {
	var loads, cycle []int
	var load int
	loadCycle := false
	cycleIndex := 1
	for i := 1; i <= times; i++ {
		if !loadCycle {
			cyclePlatform(platform)
		  logPlatform(platform, fmt.Sprintf("After %d cycles:", i), logs)
			load = getLoad(platform)
			loads = append(loads, load)
			cycle, loadCycle = checkForLoadCycle(loads)
		} else {
			load = cycle[cycleIndex]
			cycleIndex++
			if cycleIndex >= len(cycle) {
				cycleIndex = 0
			}
		}

		//log.Printf("Load @%d is %d", i, load)
	}
	return load
}

func checkForLoadCycle(loads []int) ([]int, bool) {
	tracker := make(map[int][]int)

	foundCycle := false
	for index, val := range loads {
		existing, found := tracker[val]
		if found {
			// We have seend this number before
			for _, existingIndex := range existing {
				// how many indexes between the pair
				gap := index - existingIndex

				// Compare back through the slice to ensure a cyclic pattern (3)
				for i := 0; i <= gap; i++ {
					first := existingIndex - i
					second := existingIndex - gap - i

					if first < 0 || second < 0 {
						break
					}

					foundCycle = loads[existingIndex - i] == loads[existingIndex - gap - i]
					if !foundCycle {
						break
					}
				}

				if foundCycle {
					log.Printf("Cycle Found! from %d is %d", loads, loads[existingIndex:index])
					return loads[existingIndex:index], foundCycle
				}
			}

			tracker[val] = append(tracker[val], index)
		} else {
			tracker[val] = []int{index}
		}
	}

	return []int{}, false
}

func cyclePlatform(platform [][]string) {
	rollNorthSouth(platform, 0)
	rollEastWest(platform, 0)
	rollNorthSouth(platform, len(platform) - 1)
	rollEastWest(platform, len(platform[0]) - 1)
}

func rollNorthSouth(platform [][]string, targetRow int) {
	for colIndex, symbol := range platform[targetRow] {
		rollLoc := targetRow

		currRow := targetRow
		currSymbol := symbol
		for {
			//log.Printf("Col: %d Row: %d Symbol: %s", colIndex, currRow, currSymbol)
			// Swap rolling stone to current roll location when found
			if currSymbol == "O" {
				//log.Printf("Swapping %d:%s with %d:%s", currRow, currSymbol, rollLoc, platform[rollLoc][colIndex])
				platform[currRow][colIndex] = platform[rollLoc][colIndex]
				platform[rollLoc][colIndex] = currSymbol
				rollLoc = nextIndex(rollLoc, targetRow)
			}

			// Move roll location to after the stuck rock
			if currSymbol == "#" {
				rollLoc = nextIndex(currRow, targetRow)
			}

			if testEnd(currRow, targetRow, len(platform)) {
				break
			}

			currRow = nextIndex(currRow, targetRow)
			currSymbol = platform[currRow][colIndex]
		}
	}
}

func rollEastWest(platform [][]string, targetCol int) {
	for rowIndex, symbols := range platform {
		rollLoc := targetCol

		currCol := targetCol
		currSymbol := symbols[currCol]
		for {
			//log.Printf("Col: %d Row: %d Symbol: %s", colIndex, currRow, currSymbol)
			// Swap rolling stone to current roll location when found
			if currSymbol == "O" {
				platform[rowIndex][currCol] = platform[rowIndex][rollLoc]
				platform[rowIndex][rollLoc] = currSymbol
				rollLoc = nextIndex(rollLoc, targetCol)
			}

			// Move roll location to after the stuck rock
			if currSymbol == "#" {
				rollLoc = nextIndex(currCol, targetCol)
			}

			if testEnd(currCol, targetCol, len(symbols)) {
				break
			}

			currCol = nextIndex(currCol, targetCol)
			currSymbol = platform[rowIndex][currCol]
		}
	}
}

func testEnd(currentIndex, targetIndex, rowCount int) bool {
	if targetIndex == 0 {
		return currentIndex >= rowCount - 1
	} else {
		return currentIndex <= 0
	}
}

func nextIndex(current int, target int) int {
	if target == 0 {
		return current + 1
	} else {
		return current - 1
	}
}

func countSymbols(symbols []string, target string) int {
	return strings.Count(strings.Join(symbols, ""), target)
}

func logPlatform(platform [][]string, id string, logs bool) {
	if !logs {
		return
	}
	log.Println(id)
	for _, symbols := range platform {
		log.Println(strings.Join(symbols, ""))
	}
}