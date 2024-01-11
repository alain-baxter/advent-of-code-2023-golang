package main

import (
	"alain-baxter/aoc-2023/day-10/pipe"
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
	var cmds1, cmds2 []Command
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		cmds1 = append(cmds1, parseCommand(text))
		cmds2 = append(cmds2, parseCommandPart2(text))
	}

	pipeBasedSolution(cmds1)
	projectionSolution(cmds2)
}

type Command struct {
	Direction string
	Count int
}

func parseCommand(text string) Command {
	fields := strings.Fields(text)
	direction := fields[0]
	count, _ := strconv.Atoi(fields[1])
	return Command{direction, count}
}

var directionCode = map[string]string{
	"0": "R",
	"1": "D",
	"2": "L",
	"3": "U",
}

func parseCommandPart2(text string) Command {
	fields := strings.Fields(text)
	direction := directionCode[fields[2][7:8]]
	count, _ := strconv.ParseInt(fields[2][2:7], 16, 32)
	return Command{direction, int(count)}
}

type Point struct {
	x, y int
}

var OffsetDirection = map[string]string{
	"U": "R",
	"R": "D",
	"D": "L",
	"L": "U",
}

func projectionSolution(cmds []Command) {
	current := Point{0, 0}

	verticies := []Point{current}
	perimiter := 0
	for _, cmd := range cmds {
		switch cmd.Direction {
		case "R":
			current.x += cmd.Count
		case "D":
			current.y -= cmd.Count
		case "L":
			current.x -= cmd.Count
		case "U":
			current.y += cmd.Count
		}
		verticies = append(verticies, current)
		perimiter += cmd.Count
	}
	log.Printf("Verticies: %v", verticies)

	mostArea := findAreaShoelace(verticies)

	// Since the trenches have width, the verticies are somewhere in each 1x1 square
	// and so not all the area is calculated by Shoelace. If you think of the vertex being
	// in the top left corner, then all the right and all the bottom border 1x1 squares 
	// are missed. These are accounted for by half of the perimiter + 1
	area := mostArea + (perimiter / 2) + 1
	log.Printf("Part 2 Final volume: %d", area)
}

func findAreaShoelace(verticies []Point) int {
	area := 0

	for i := 0; i < len(verticies) - 1; i++ {
		pointA, pointB := verticies[i], verticies[i + 1]
		area += (pointA.y + pointB.y) * (pointA.x - pointB.x)
	}

	return abs(area / 2)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

/*
  This is the solution I used for part 1. I recognized this was basically the same
	as day-10 but the pipe was described by perimiter with the instructions instead of
	given outright. So this solution created the pipe diagram from the instructions
	then uses the same intersection code to get inside vs outside for volume.

	This did not scale well to part 2 sized diagrams. It couldn't even do the example
	diagram without crashing OOM.
*/
func pipeBasedSolution(cmds []Command) {
	var area [][]pipe.Pipe
	startShape := "S"
	start := pipe.Pipe{Row: 0, Col: 0, Shape: startShape, Distance: 1}
	area = append(area, []pipe.Pipe{start})

	current := start
	volume := 0
	for _, cmd := range cmds {
		current, area = addToArea(cmd, current, area)
		volume += cmd.Count
	}
	area = fixLastCorner(current, area)

	areaAsMap := convertAreaToMap(area)
	for i, row := range area {
		for j, section := range row {
			if section.InCycle() {
				continue
			}

			section.Row = i
			section.Col = j
			intersections := section.GetIntersections(areaAsMap, false)
			if intersections%2 != 0 {
				section.Shape = "#"
				volume++
			}
			area[i][j] = section
		}
	}
	logArea(area, "Counted Area:")
	log.Printf("Part 1 Final volume: %d", volume)
}

func fixLastCorner(current pipe.Pipe, area [][]pipe.Pipe) [][]pipe.Pipe {
	up, down, left, right := "", "", "", ""
	if current.Row - 1 > 0 {
		up = area[current.Row - 1][current.Col].Shape
	}
	if current.Row + 1 < len(area) {
		down = area[current.Row + 1][current.Col].Shape
	}
	if current.Col - 1 > 0 {
		left = area[current.Row][current.Col - 1].Shape
	}
	if current.Col + 1 < len(area[0]) {
		right = area[current.Row][current.Col + 1].Shape
	}

	if up == "|" && right == "-" {
		current.Shape = "L"
	} else if up == "|" && left == "-" {
		current.Shape = "J"
	} else if down == "|" && right == "-" {
		current.Shape = "F"
	} else if down == "|" && left == "-" {
		current.Shape = "7"
	}

	area[current.Row][current.Col] = current	
	return area
}

func addToArea(cmd Command, current pipe.Pipe, area [][]pipe.Pipe) (pipe.Pipe, [][]pipe.Pipe) {
	switch cmd.Direction {
	case "R":
		return goRight(cmd.Count, current, area)
	case "D":
		return goDown(cmd.Count, current, area)
	case "L":
		return goLeft(cmd.Count, current, area)
	case "U":
		return goUp(cmd.Count, current, area)
	}
	return current, area
}

func goRight(count int, current pipe.Pipe, area [][]pipe.Pipe) (pipe.Pipe, [][]pipe.Pipe) {
	area, current = padArea(0, count, current, area)
  colStart := current.Col + 1
	colEnd := colStart + count

	// Ensure the correct direction on current
	if current.Shape == "7" {
		current.Shape = "F"
		area[current.Row][current.Col] = current
	} else if current.Shape == "J" {
		current.Shape = "L"
		area[current.Row][current.Col] = current
	}

	var r pipe.Pipe
	for col := colStart; col < colEnd; col++ {
		shape := "-"
		if col == colEnd - 1 {
			shape = "7"
		}
		r = pipe.Pipe{Row: current.Row, Col: col, Shape: shape, Distance: 1}
		area[current.Row][col] = r
	}
	//logArea(area, "New Area:")
	return r, area
}

func goDown(count int, current pipe.Pipe, area [][]pipe.Pipe) (pipe.Pipe, [][]pipe.Pipe) {
	area, current = padArea(count, 0, current, area)
  rowStart := current.Row + 1
	rowEnd := rowStart + count

	// Ensure the correct direction on current
	if current.Shape == "J" {
		current.Shape = "7"
		area[current.Row][current.Col] = current
	} else if current.Shape == "L" {
		current.Shape = "F"
		area[current.Row][current.Col] = current
	}

	var r pipe.Pipe
	for row := rowStart; row < rowEnd; row++ {
		shape := "|"
		if row == rowEnd - 1 {
			shape = "J"
		}
		r = pipe.Pipe{Row: row, Col: current.Col, Shape: shape, Distance: 1}
		area[row][current.Col] = r
	}
	//logArea(area, "New Area:")
	return r, area
}

func goLeft(count int, current pipe.Pipe, area [][]pipe.Pipe) (pipe.Pipe, [][]pipe.Pipe) {
	area, current = padArea(0, -count, current, area)
	colStart := current.Col - 1
	colEnd := colStart - count

	// Ensure the correct direction on current
	if current.Shape == "L" {
		current.Shape = "J"
		area[current.Row][current.Col] = current
	} else if current.Shape == "F" {
		current.Shape = "7"
		area[current.Row][current.Col] = current
	}

	var r pipe.Pipe
	for col := colStart; col > colEnd; col-- {
		shape := "-"
		if col == colEnd + 1 {
			shape = "L"
		}
		r = pipe.Pipe{Row: current.Row, Col: col, Shape: shape, Distance: 1}
		area[current.Row][col] = r
	}
	//logArea(area, "New Area:")
	return r, area
}

func goUp(count int, current pipe.Pipe, area [][]pipe.Pipe) (pipe.Pipe, [][]pipe.Pipe) {
	area, current = padArea(-count, 0, current, area)
	rowStart := current.Row - 1
	rowEnd := rowStart - count

	// Ensure the correct direction on current
	if current.Shape == "7" {
		current.Shape = "J"
		area[current.Row][current.Col] = current
	} else if current.Shape == "F" {
		current.Shape = "L"
		area[current.Row][current.Col] = current
	}

	var r pipe.Pipe
	for row := rowStart; row > rowEnd; row-- {
		shape := "|"
		if row == rowEnd + 1 {
			shape = "F"
		}
		r = pipe.Pipe{Row: row, Col: current.Col, Shape: shape, Distance: 1}
		area[row][current.Col] = r
	}
	//logArea(area, "New Area:")
	return r, area
}

func padArea(rowOffset, colOffset int, currentLoc pipe.Pipe, area [][]pipe.Pipe) ([][]pipe.Pipe, pipe.Pipe) {
	currentRows := len(area)
	currentCols := len(area[0])
	
	colLeftPad := 0
	colRightPad := 0
	if colOffset < 0 {
		colLeftPad = currentLoc.Col + colOffset
		if colLeftPad > 0 {
			colLeftPad = 0
		} else {
			colLeftPad = -1 * colLeftPad
		}
		// Move current location
		currentLoc.Col += colLeftPad
	} else if colOffset > 0 && currentLoc.Col + colOffset >= currentCols {
		colRightPad = currentLoc.Col + colOffset - (currentCols - 1)
	}
	fullColLength := colLeftPad + currentCols + colRightPad

	rowUpPad := 0
	rowDownPad := 0
	if rowOffset < 0 {
		rowUpPad = currentLoc.Row + rowOffset
		if rowUpPad > 0 {
			rowUpPad = 0
		} else {
			rowUpPad = -1 * rowUpPad
		}
		currentLoc.Row += rowUpPad
	} else if rowOffset > 0 && currentLoc.Row + rowOffset >= currentRows {
		rowDownPad = currentLoc.Row + rowOffset - (currentRows - 1)
	}
	totalRows := rowUpPad + currentRows + rowDownPad

	for i := 0; i < totalRows; i++ {
		// brand new front padded rows
		if i < rowUpPad {
			newRow := make([]pipe.Pipe, fullColLength)
			area = append([][]pipe.Pipe{newRow}, area...)
		
		// brand new back padded rows
		} else if i > rowUpPad + currentRows - 1 {
			area = append(area, make([]pipe.Pipe, fullColLength))
		
		// Extend existing rows with padding
		} else {
			row := append(make([]pipe.Pipe, colLeftPad), area[i]...)
			row = append(row, make([]pipe.Pipe, colRightPad)...)
			area[i] = row
		}
	}
	return area, currentLoc
}

func convertAreaToMap(area [][]pipe.Pipe) map[string]pipe.Pipe {
	r := make(map[string]pipe.Pipe)
	for i, row := range area {
		for j, pipe := range row {
			pipe.Row = i
			pipe.Col = j
			r[pipe.GetName()] = pipe
		}
	}
	return r
}

func logArea(area [][]pipe.Pipe, id string) {
	log.Println(id)
	for _, row := range area {
		var buffer string
		for _, pipe := range row {
			if len(pipe.Shape) == 0 {
				buffer += "."
			} else {
				buffer += pipe.Shape
			}
		}
		log.Println(buffer)
	}
}