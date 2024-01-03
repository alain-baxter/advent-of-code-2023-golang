package main

import (
	"alain-baxter/aoc-2023/day-10/pipe"
	"bufio"
	"log"
	"maps"
	"os"
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

	pipes := make(map[string]pipe.Pipe)
	row := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		maps.Copy(pipes, pipe.GetPipes(text, row))
		row++
	}
	log.Println(pipes)

	var start pipe.Pipe
	for _, v := range pipes {
		if v.Shape == "S" {
			start = v
			break
		}
	}
	log.Printf("Starting pipe: %v",start)

	branch1, branch2 := pipe.FollowStart(pipes, start)
	prev1, prev2 := start, start
	log.Printf("Branch1 start: %v", branch1)
	log.Printf("Branch2 start: %v", branch2)
	var end1, end2 bool = false, false

	for {
		branch1, prev1, end1 = pipe.FollowPipe(pipes, branch1, prev1.GetName())
		branch2, prev2, end2 = pipe.FollowPipe(pipes, branch2, prev2.GetName())
		log.Printf("branch1 %v", branch1)
		log.Printf("branch2 %v", branch2)

		if end1 {
			break
		}
		if end2 {
			break
		}
	}

	end := branch1
	if branch2.Distance > branch1.Distance {
		end = branch2
	}

	log.Printf("Steps to furthest point from start: %v", end.Distance)

	inside := 0
	for _, v := range pipes {
		if v.InCycle() {
			continue
		}

		intersections := v.GetIntersections(pipes)
		if intersections%2 != 0 {
			inside++
		}
	}

	log.Printf("Number of inside points: %d", inside)
}