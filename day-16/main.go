package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const NORTH byte = 0x01
const EAST byte = 0x02
const SOUTH byte = 0x04
const WEST byte = 0x08

const R_MIRROR string = "\\"
const L_MIRROR string = "/"
const V_SPLITTER string = "|"
const H_SPLITTER string = "-"
const EMPTY string = "."

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
	var contraption [][]string
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		contraption = append(contraption, strings.Split(text, ""))
	}

	log.Printf("Contraption: %s", contraption)

	energized := followLight(Beam{0, 0, EAST}, contraption, true)

	log.Printf("Part 1 Energized count: %d", energized)

	maxEnergized := 0
	var maxBeam Beam
	// Check Top and Bottom edges
	for _, x := range []int{0, len(contraption) - 1}  {
		direction := NORTH
		if x == 0 {
			direction = SOUTH
		}
		for y := 0; y < len(contraption[0]); y++ {
			energized := followLight(Beam{x, y, direction}, contraption, false)
			if energized > maxEnergized {
				maxEnergized = energized
				maxBeam = Beam{x, y, direction}
			}
		}
	}

	// Check Left and Right edges
	for _, y := range []int{0, len(contraption[0]) - 1}  {
		direction := WEST
		if y == 0 {
			direction = EAST
		}
		for x := 0; x < len(contraption); x++ {
			energized := followLight(Beam{x, y, direction}, contraption, false)
			if energized > maxEnergized {
				maxEnergized = energized
				maxBeam = Beam{x, y, direction}
			}
		}
	}

	log.Printf("Part 2 Max energized: %d with beam %v", maxEnergized, maxBeam)
}

func decodeEnergized(row []byte) string {
	var r string
	for _, val := range row {
		if val > 0 {
			r = r + "#"
		} else {
			r = r + "."
		}
	}
	return r
}

func followLight(starting Beam, contraption [][]string, doLog bool) int {
	beams := []Beam{starting}

	energy := make([][]byte, len(contraption))
	for i := range energy {
			energy[i] = make([]byte, len(contraption[0]))
	}
	energy[starting.X][starting.Y] += starting.Direction

	//log.Printf("Starting Beam: %v", beams)
	for {
		beams = nextBeams(beams, contraption, &energy)
		//log.Printf("Next Beam: %v", beams)

		if len(beams) == 0 {
			break
		}
	}

	if doLog {
		log.Printf("Enery levels: %d", energy)
		log.Println("Energized Tiles")
	}

	energized := 0
	for _, row := range energy {
		if doLog {
			log.Println(decodeEnergized(row))
		}
		for _, level := range row {
			if level > 0 {
				energized++
			}
		}
	}

	//log.Printf("Starting %v energized %d", starting, energized)
	return energized
}

func nextBeams(beams []Beam, contraption [][]string, energy *[][]byte) []Beam {
	var buffer, r []Beam

	for _, beam := range beams {
		switch symbol := contraption[beam.X][beam.Y]; symbol {
		case R_MIRROR:
			buffer = append(buffer, rightMirror(beam))
		case L_MIRROR:
			buffer = append(buffer, leftMirror(beam))
		case H_SPLITTER:
			buffer = append(buffer, horizontalSplitter(beam)...)
		case V_SPLITTER:
			buffer = append(buffer, verticalSplitter(beam)...)
		default:
			buffer = append(buffer, sameDirection(beam))
		}
	}

	xMax := len(contraption)
	yMax := len(contraption[0])
	for _, beam := range buffer {
		if beam.X < xMax && beam.Y < yMax && beam.X >= 0 && beam.Y >= 0 {
			// Check is we have entered from this direction before
			currentEnergy := (*energy)[beam.X][beam.Y]
			if currentEnergy & beam.Direction != beam.Direction {
				r = append(r, beam)
			  (*energy)[beam.X][beam.Y] += beam.Direction
			}
		}
	}
	return r
}

func rightMirror(beam Beam) Beam {
	switch beam.Direction {
	case NORTH:
		beam.Direction = WEST
		beam.Y--
	case EAST:
		beam.Direction = SOUTH
		beam.X++
	case SOUTH:
		beam.Direction = EAST
		beam.Y++
	case WEST:
		beam.Direction = NORTH
		beam.X--
	}
	return beam
}

func leftMirror(beam Beam) Beam {
	switch beam.Direction {
	case NORTH:
		beam.Direction = EAST
		beam.Y++
	case EAST:
		beam.Direction = NORTH
		beam.X--
	case SOUTH:
		beam.Direction = WEST
		beam.Y--
	case WEST:
		beam.Direction = SOUTH
		beam.X++
	}
	return beam
}

func horizontalSplitter(beam Beam) []Beam {
	switch beam.Direction {
	case EAST:
		return []Beam{sameDirection(beam)}
	case WEST:
		return []Beam{sameDirection(beam)}
	}
	split := Beam{beam.X, beam.Y - 1, WEST}
	beam.Direction = EAST
	beam.Y++
	return []Beam{beam, split}
}

func verticalSplitter(beam Beam) []Beam {
	switch beam.Direction {
	case NORTH:
		return []Beam{sameDirection(beam)}
	case SOUTH:
		return []Beam{sameDirection(beam)}
	}
	split := Beam{beam.X - 1, beam.Y, NORTH}
	beam.Direction = SOUTH
	beam.X++
	return []Beam{beam, split}
}

func sameDirection(beam Beam) Beam {
	switch beam.Direction {
	case NORTH:
		beam.X--
	case EAST:
		beam.Y++
	case SOUTH:
		beam.X++
	case WEST:
		beam.Y--
	}
	return beam
}

type Beam struct {
	X, Y int
	Direction byte
}