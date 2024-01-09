package main

import (
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
	boxes := make([][]Lens, 256)
	var cmds []string
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		cmds = strings.Split(text, ",")
	}
	log.Printf("Cmds: %s", cmds)

	hashSum := 0
	for _, cmd := range cmds {
		hashSum += hash(cmd)
	}
	log.Printf("Hasing Sum: %d", hashSum)

	for _, cmd := range cmds {
		lens := parseLens(cmd)
		hash := hash(lens.Label)
		box := boxes[hash]
		if lens.FocalLength > 0 {
			index, found := checkForLens(box, lens)
			if found {
				box[index] = lens
			} else {
				box = append(box, lens)
			}
			boxes[hash] = box
		} else {
			index, found := checkForLens(box, lens)
			if found {
				box := append(box[:index], box[index+1:]...)
				boxes[hash] = box
			}
		}
	}
	log.Printf("Final boxes: %v", boxes)
	
	focusingPower := 0
	for index, box := range boxes {
		for slot, lens := range box {
			focusingPower += (index + 1) * (slot + 1) * lens.FocalLength
		}
	}
	log.Printf("Focusing Power: %d", focusingPower)
}

func hash(s string) int {
	hash := 0
	for _, r := range s {
		hash += int(r)
		hash = (hash * 17) % 256
	}
	return hash
}

func checkForLens(box []Lens, lens Lens) (int, bool) {
	for index, val := range box {
		if val.Label == lens.Label {
			return index, true
		}
	}
	return -1, false
}

func parseLens(cmd string) Lens {
	if strings.Contains(cmd, "=") {
		parts := strings.Split(cmd, "=")
		focal, _ := strconv.Atoi(parts[1])
		return Lens{parts[0], focal}
	} else {
		return Lens{cmd[:len(cmd)-1], -1}
	}
}

type Lens struct {
	Label string
	FocalLength int
}