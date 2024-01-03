package main

import (
	"alain-baxter/aoc-2023/day-8/node"
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

	var startingNodes []node.Node
	nodes := make(map[string]node.Node)
	var directions []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)

		if len(directions) == 0 {
			directions = strings.Split(text, "")
			continue
		}

		if len(text) == 0 {
			continue
		}

		node := node.NewNode(text)
		if strings.LastIndex(node.Label, "A") == 2 {
			startingNodes = append(startingNodes, node)
		}
		nodes[node.Label] = node
	}

	log.Printf("Directions: %s", directions)
	log.Printf("Node Map: %v", nodes)
	log.Printf("Starting Nodes: %v", startingNodes)

	cycles := make(map[string]int)
	for _, node := range startingNodes {
		cycles[node.Label] = getCycle(node, nodes, directions)
	}

	log.Printf("Cycles: %v", cycles)
	
	var cyclesList []int
	for _, cycle := range cycles {
		cyclesList = append(cyclesList, cycle)
	}
	log.Printf("LCM of %v is %d", cycles, LCM(cyclesList[0], cyclesList[1], cyclesList[2:]...))
}

func getCycle(starting node.Node, nodes map[string]node.Node, directions []string) int {
	label := starting.Label
	var index, cycle int = 0, 0
	currNode := starting

	for {
		if index >= len(directions) {
			index = 0
		}

		if directions[index] == "L" {
			currNode = nodes[currNode.Left]
		} else {
			currNode = nodes[currNode.Right]
		}

		index++
		cycle++
		if strings.LastIndex(currNode.Label, "Z") == 2 {
			log.Printf("%s reached %s at index %d with cycle %d", label, currNode.Label, index, cycle)
			return cycle
		}
	}
}

func GCDEuclidean(a, b int) int {
	for a != b {
		if a > b {
			a -= b
		} else {
			b -= a
		}
	}

	return a
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
      for b != 0 {
              t := b
              b = a % b
              a = t
      }
      return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
      result := a * b / GCD(a, b)

      for i := 0; i < len(integers); i++ {
              result = LCM(result, integers[i])
      }

      return result
}
