package main

import (
	"alain-baxter/aoc-2023/utils"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	lines := utils.LoadFile(os.Args)
	steps := getSteps(os.Args)
	garden, start := parseGarden(lines)
	if steps > len(garden)/2 {
		doStepsGeometrically(&garden, start, steps)
	} else {
		countReachable := doSteps(&garden, start, steps)
		log.Println("(Part 1) Reachable in", steps, "steps:", countReachable)
	}
}

type Tile struct {
	x, y       int
	gardenPlot bool
}

func getSteps(args []string) int {
	if len(args) > 2 {
		steps, err := strconv.Atoi(args[2])
		if err != nil {
			log.Println("Invalid step", args[2])
		} else {
			return steps
		}
	}

	return 64
}

func parseGarden(lines []string) ([][]Tile, Tile) {
	// Initialize the garden size
	tiles := make([][]Tile, len(lines))
	for i := 0; i < len(lines); i++ {
		tiles[i] = make([]Tile, len(lines[0]))
	}

	var start Tile
	for x, line := range lines {
		for y, symbol := range strings.Split(line, "") {
			plot := symbol == "." || symbol == "S"
			tile := Tile{x, y, plot}
			tiles[x][y] = tile

			if symbol == "S" {
				start = tile
			}
		}
	}

	return tiles, start
}

func doSteps(garden *[][]Tile, start Tile, steps int) int {
	active := []Tile{start}
	evenTiles := []Tile{start}
	oddTiles := []Tile{}
	for step := 1; step <= steps; step++ {
		newActive := []Tile{}

		var seen *[]Tile
		if step%2 == 0 {
			seen = &evenTiles
		} else {
			seen = &oddTiles
		}

		for _, tile := range active {
			newActive = append(newActive, getNext(tile, garden, seen, step)...)
		}

		active = newActive
	}

	var count int
	if steps%2 == 0 {
		count = len(evenTiles)
	} else {
		count = len(oddTiles)
	}
	return count
}

func getNext(tile Tile, garden *[][]Tile, seen *[]Tile, step int) []Tile {
	possible := []Tile{}

	down := tile
	down.x++
	possible = append(possible, down)
	up := tile
	up.x--
	possible = append(possible, up)
	right := tile
	right.y++
	possible = append(possible, right)
	left := tile
	left.y--
	possible = append(possible, left)

	r := []Tile{}
	maxX := len(*garden)
	maxY := len((*garden)[0])
	for _, maybe := range possible {
		check := (*garden)[utils.ModProperly(maybe.x, maxX)][utils.ModProperly(maybe.y, maxY)]
		// If a rock, skip
		if !check.gardenPlot {
			continue
		}

		if slices.Contains(*seen, maybe) {
			continue
		}

		*seen = append(*seen, maybe)
		r = append(r, maybe)
	}

	return r
}

func doStepsGeometrically(garden *[][]Tile, start Tile, steps int) {
	evenCount, oddCount, evenCorners, oddCorners := countGarden(garden, start)
	log.Println(evenCount, oddCount, evenCorners, oddCorners)

	rowSize := len((*garden)[0])
	// The number of repetitions of the garden for the given steps
	n := (steps - start.x) / rowSize
	log.Println("Row Size:", rowSize, "Times:", n)
	if steps%2 == 0 {
		log.Fatal("Only created count for odd steps geometrically since that was the problem")
	} else {
		// Laying out an example. Since rowSize is odd the garden will switch between even and odd parity nodes being
		// activated. This leads to (n+1)^2 odd parity full sections, and n^2 even parity full sections. Since the axis
		// are all "." and there is a diamond swath of "." we also reach the edges perfectly. So we can add and subtract
		// some odd and even parity corners to smooth out the edge of the multi-region diamond and get the proper count.
		count := ((n + 1) * (n + 1) * oddCount) + (n * n * evenCount) - ((n + 1) * oddCorners) + (n * evenCorners)
		log.Println("(Part 2) count of possible steps:", count)
	}
}

func countGarden(garden *[][]Tile, start Tile) (int, int, int, int) {
	log.Println("Start:", start)
	var evenCount, oddCount, evenCorners, oddCorners int = 0, 0, 0, 0
	for x, row := range *garden {
		for y, tile := range row {
			// Skip rocks
			if !tile.gardenPlot {
				continue
			}

			// Tile is surrounded by rocks
			if len(getNext(tile, garden, &[]Tile{}, 0)) == 0 {
				continue
			}

			// Move to "start" being in @ (0, 0)
			translatedX := x - start.x
			translatedY := y - start.y
			distance := utils.Abs(translatedX) + utils.Abs(translatedY)

			// Even point
			even := false
			if distance%2 == 0 {
				evenCount++
				even = true
				// Odd Point
			} else {
				oddCount++
			}

			corner := distance > start.x
			if corner && even {
				evenCorners++
			} else if corner && !even {
				oddCorners++
			}
		}
	}

	return evenCount, oddCount, evenCorners, oddCorners
}
