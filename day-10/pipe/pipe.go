package pipe

import (
	"fmt"
	"log"
	"slices"
)

type Pipe struct {
	Row, Col int
	Shape string
	Distance int
}

func (p Pipe) GetName() string {
	return fmt.Sprint(p.Row) + ":" + fmt.Sprint(p.Col)
}

func (p Pipe) offsetName(rowOff, colOff int) string {
	return fmt.Sprint(p.Row + rowOff) + ":" + fmt.Sprint(p.Col + colOff)
}

func GetPipes(s string, row int) map[string]Pipe {
	ps := make(map[string]Pipe)
	for i := range s {
		shape := string(s[i])
		pipe := Pipe{Shape: shape, Row: row, Col: i}
		ps[pipe.GetName()] = pipe
	}
	return ps
}

func FindStart(pipes map[string]Pipe) Pipe {
	for _, v := range pipes {
		if (v.Shape == "S") {
			v.Distance = 0
			return v
		}
	}
	return Pipe{}
}

func FollowPipe(pipes map[string]Pipe, pipe Pipe, prev string) (Pipe, Pipe, bool) {
	var rowOff1, colOff1, rowOff2, colOff2 int = 0, 0, 0, 0
	switch {
	case pipe.Shape == "L":
		rowOff1 = -1
		colOff2 = 1
	case pipe.Shape == "7":
		colOff1 = -1
		rowOff2 = 1
	case pipe.Shape == "J":
		rowOff1 = -1
		colOff2 = -1
	case pipe.Shape == "F":
		colOff1 = 1
		rowOff2 = 1
	case pipe.Shape == "|":
		rowOff1 = -1
		rowOff2 = 1
	case pipe.Shape == "-":
		colOff1 = -1
		colOff2 = 1
	}

	name := pipe.offsetName(rowOff1, colOff1)
	if name == prev {
		name = pipe.offsetName(rowOff2, colOff2)
	}
	next := pipes[name]

	if next.Distance == 0 {
		next.Distance = pipe.Distance + 1
		pipes[name] = next
		return next, pipe, false
	} else {
		pipes[name] = next
		return next, pipe, true
	}
}

func FollowStart(pipes map[string]Pipe, pipe Pipe) (Pipe, Pipe) {
	var possible []Pipe

	n, nok := pipes[pipe.offsetName(-1, 0)]
	if nok && slices.Contains([]string{"F", "7", "|"}, n.Shape) {
		n.Distance = 1
		possible = append(possible, n)
	}
	s, sok := pipes[pipe.offsetName(1, 0)]
	if sok && slices.Contains([]string{"L", "J", "|"}, s.Shape) {
		s.Distance = 1
		possible = append(possible, s)
	}
	e, eok := pipes[pipe.offsetName(0, 1)]
	if eok && slices.Contains([]string{"7", "J", "-"}, e.Shape) {
		e.Distance = 1
		possible = append(possible, e)
	}
	w, wok := pipes[pipe.offsetName(0, -1)]
	if wok && slices.Contains([]string{"L", "F", "-"}, w.Shape) {
		w.Distance = 1
		possible = append(possible, w)
	}

	if len(possible) != 2 {
		log.Fatalf("Too many outlets from S %v", possible)
	}

	pipes[possible[0].GetName()] = possible[0]
	pipes[possible[1].GetName()] = possible[1]
	return possible[0], possible[1]
}

func (p Pipe) GetIntersections(pipes map[string]Pipe, doLog bool) int {
	intersections := 0
	current := p
	ok := true
	lastShape := "."

	for {
		name := current.offsetName(1, 0)
		current, ok = pipes[name]
		if ok {
			if current.InCycle() && current.Shape != "|" {
				// Consider this a single intersection similar to /
				if lastShape == "F" && current.Shape == "J" {
					lastShape = current.Shape
				// Consider this a single intersection similar to \
				} else if lastShape == "7" && current.Shape == "L" {
					lastShape = current.Shape
				} else {
					lastShape = current.Shape
					intersections++
				}
			}
		} else {
			break
		}
	}

	if doLog {
	  log.Printf("Intercetions from %v: %d", p, intersections)
	}
	return intersections
}

func (p Pipe) InCycle() bool {
	return p.Shape == "S" || p.Distance != 0
}