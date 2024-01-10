package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const NORTH string = "N"
const EAST string = "E"
const SOUTH string = "S"
const WEST string = "W"

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
	var heatMap [][]int
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		heatMap = append(heatMap, toInts(strings.Split(text, "")))
	}
	log.Printf("Map: %d", heatMap)

	bestPath, bestHeatLoss := followPaths([]Path{{0, 0, 0, SOUTH},{0, 0, 0, EAST}}, Crucible{0, 3}, &heatMap)
	log.Printf("Best path normal crucible %v with heat loss: %d", bestPath, bestHeatLoss)

	bestPath, bestHeatLoss = followPaths([]Path{{0, 0, 0, SOUTH},{0, 0, 0, EAST}}, Crucible{4, 10}, &heatMap)
	log.Printf("Best path ultra crucible %v with heat loss: %d", bestPath, bestHeatLoss)
}

func toInts(strs []string) []int {
	ints := make([]int, len(strs))
	for i, item := range strs {
			val, _ := strconv.Atoi(item)
			ints[i] = val
	}
	return ints
}

func followPaths(starting []Path, crucible Crucible, heatMap *[][]int) (Path, int) {
	tracker := make(map[Path]int)
	for _, path := range starting {
		tracker[path] = 0
	}

	bestHeatLoss := math.MaxInt32
	var bestPath Path
	for {
		currentTracker := make(map[Path]int)
		for path, heatLoss := range tracker {
			nextPaths := getNextPaths(path, heatMap, crucible.minSteps)
			for _, nextPath := range nextPaths {
				// Toss any paths with more than maximum steps in a directions
				if nextPath.Steps > crucible.maxSteps {
					continue
				}

				newHeatLoss := (*heatMap)[nextPath.X][nextPath.Y] + heatLoss

				// Check if we made it to the end
				if nextPath.X == len(*heatMap) - 1 && nextPath.Y == len((*heatMap)[0]) - 1 {
					if newHeatLoss < bestHeatLoss && nextPath.Steps >= crucible.minSteps {
						bestHeatLoss = newHeatLoss
						bestPath = nextPath
					}
				}

				// Check within the current tracker for matching path
				existing, foundCurr := currentTracker[nextPath]
				if foundCurr {
					currentTracker[nextPath] = int(math.Min(float64(existing), float64(newHeatLoss)))
				}

				// Check within the previous tracker for matching path
				existing, foundPrev := tracker[nextPath]
				if foundPrev {
					currentTracker[nextPath] = int(math.Min(float64(existing), float64(newHeatLoss)))
				}

				if !foundCurr && !foundPrev && newHeatLoss < bestHeatLoss {
					currentTracker[nextPath] = newHeatLoss
				}
			}
		}

		tracker = currentTracker

		if len(tracker) == 0 {
			break
		}
	}
	return bestPath, bestHeatLoss
}

func getNextPaths(path Path, heatMap *[][]int, minSteps int) []Path {
	var buffer, r []Path

	switch path.Direction {
	case NORTH:
		buffer = append(buffer, path.goNorth())
		if path.Steps >= minSteps {
			buffer = append(buffer, path.goEast(), path.goWest())
		}

	case EAST:
		buffer = append(buffer, path.goEast())
		if path.Steps >= minSteps {
			buffer = append(buffer, path.goNorth(), path.goSouth())	
		}

	case SOUTH:
		buffer = append(buffer, path.goSouth())
		if path.Steps >= minSteps {
			buffer = append(buffer, path.goEast(), path.goWest())
		}

	case WEST:
		buffer = append(buffer, path.goWest())
		if path.Steps >= minSteps {
			buffer = append(buffer, path.goNorth(), path.goSouth())
		}
	}

	xMax := len(*heatMap)
	yMax := len((*heatMap)[0])
	for _, next := range buffer {
		if next.X < xMax && next.Y < yMax && next.X >= 0 && next.Y >= 0 {
			r = append(r, next)
		}
	}
	return r
}

func (p Path) goNorth() Path {
	steps := p.Steps + 1
	if p.Direction != NORTH {
		steps = 1
	}
	return Path{p.X - 1, p.Y, steps, NORTH}
}

func (p Path) goEast() Path {
	steps := p.Steps + 1
	if p.Direction != EAST {
		steps = 1
	}
	return Path{p.X, p.Y + 1, steps, EAST}
}

func (p Path) goSouth() Path {
	steps := p.Steps + 1
	if p.Direction != SOUTH {
		steps = 1
	}
	return Path{p.X + 1, p.Y, steps, SOUTH}
}

func (p Path) goWest() Path {
	steps := p.Steps + 1
	if p.Direction != WEST {
		steps = 1
	}
	return Path{p.X, p.Y - 1, steps, WEST}
}

type Path struct {
	X, Y, Steps int
	Direction string
}

type Crucible struct {
	minSteps, maxSteps int
}