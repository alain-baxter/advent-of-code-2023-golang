package main

import (
	"alain-baxter/aoc-2023/day-19/part"
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

	scanner := bufio.NewScanner(file)
	workflows := make(map[string]part.Workflow)
	workflows["A"] = part.Workflow{Label: "A"}
	workflows["R"] = part.Workflow{Label: "R"}
	var parts []part.Part
	parsePhase := 0
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		if len(strings.TrimSpace(text)) == 0 {
			parsePhase = 1
			continue
		}

		if parsePhase == 0 {
			workflow := part.ParseWorkflow(text)
		  workflows[workflow.Label] = workflow
		} else {
			parts = append(parts, part.ParsePart(text))
		}
	}

	sumAccepted := 0
	for _, p := range parts {
		sumAccepted += part.CheckPart(p, workflows)
	}
	log.Printf("(part 1) Sum of Accepted Parts: %d", sumAccepted)
	log.Printf("(part 2) %d distinct combinations of ratings that will be accepted.", part.PathsToAccepted(workflows))
}