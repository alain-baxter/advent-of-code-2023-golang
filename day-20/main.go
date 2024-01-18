package main

import (
	"alain-baxter/aoc-2023/utils"
	"log"
	"os"
	"strings"
)

func main() {
	lines := utils.LoadFile(os.Args)

	broadcast := parseInput(lines)
	queue := MessageQueue{entries: []Message{}, lowCount: 0, highCount: 0, rxPrecursorCount: 0, rxPrecursorCycles: make(map[Module]int)}
	output := queue.execute(broadcast, 1000)
	log.Println("(Part 1) After 1000 button pushes:", output)

	broadcast = parseInput(lines)
	rxIter := queue.toRX(broadcast)
	log.Println("(Part 2) First iteration that rx got a low:", rxIter)
}

func parseInput(lines []string) Broadcast {
	modules := make(map[string]Module)
	var broadcast Broadcast

	// Initialize Modules Objects
	for _, line := range lines {
		prefix := string(line[0])
		label := line[:strings.Index(line, " ")]

		if prefix == "%" {
			flipflop := FlipFlop{false, []Module{}}
			modules[label[1:]] = &flipflop
		} else if prefix == "&" {
			conjunction := Conjunction{make(map[Module]bool), []Module{}}
			modules[label[1:]] = &conjunction
		} else if label == "broadcaster" {
			broadcast = Broadcast{[]Module{}}
			modules[label] = &broadcast
		} else {
			null := Null{}
			modules[label] = &null
		}
	}

	// Setup input/output to modules
	for _, line := range lines {
		parts := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(line, "&", ""), "%", ""), "->", ""), ",", ""))
		label := parts[0]
		module := modules[label]

		for i := 1; i < len(parts); i++ {
			destination, ok := modules[parts[i]]
			if !ok {
				destination = &Null{label: parts[i]}
			}
			module.AddOutput(destination)
			destination.AddInput(module)
		}
	}

	log.Println(modules)

	return broadcast
}
