package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
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
	row := 1
	colExpansion := make(map[int]bool)
	expantionFactor := 1000000
	var stars []Star
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		rowStars := scanRow(text, row, colExpansion)

		if len(rowStars) == 0 {
			// Expanding Rows in place
			row += expantionFactor - 1
		} else {
			stars = append(stars, rowStars...)
		}
		
		row++
	}

	expandColumns(stars, colExpansion, expantionFactor)

	log.Printf("Col & Row Expanded Stars: %v", stars)

	orthogonalSum := 0
	for i, star1 := range stars {
		for _, star2 := range stars[i+1:] {
			orthogonalSum += getOrthogonalDistance(star1, star2)
		}
	}

	log.Printf("Orthogonal sum: %d", orthogonalSum)
}

func scanRow(s string, row int, colExpansion map[int]bool) []Star {
	var stars []Star
	for i := range s {
		val := string(s[i])
		col := i + 1
		if val == "#" {
			colExpansion[col] = false
			stars = append(stars, Star{row: row, col: col})
		} else {
			_, found := colExpansion[col]
			if !found {
				colExpansion[col] = true
			}
		}
	}
	return stars
}

func expandColumns(stars []Star, colExpansion map[int]bool, expantionFactor int) {
	var toBeExpanded []int
	for col, expand := range colExpansion {
		if expand {
			toBeExpanded = append(toBeExpanded, col)
		}
	}

	// expand from highest to lowest to ensure consistent answer
	sort.Slice(toBeExpanded, func(i, j int) bool { return toBeExpanded[i] > toBeExpanded[j]})
	for _, col := range toBeExpanded {
		for i, star := range stars {
			if star.col > col {
				star.col = star.col + expantionFactor - 1
				stars[i] = star
			}
		}
	}
}

func getOrthogonalDistance(star1, star2 Star) int {
	distance := 0

	distance += int(math.Abs(float64(star1.row) - float64(star2.row)))
	distance += int(math.Abs(float64(star1.col) - float64(star2.col)))

	//log.Printf("Distance from %v to %v is %d", star1, star2, distance)
	return distance
}

type Star struct {
	row, col int
}