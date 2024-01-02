package main

import (
	"bufio"
	"log"
	"os"
	"sort"
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

	var hands []Hand
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		fields := strings.Fields(text)
		cards := fields[0]
		bid, _ := strconv.Atoi(fields[1])
		hands = append(hands, newHand(cards, bid))
	}

	sort.Slice(hands, func(i, j int) bool {
		// hand is "less" when is has a high number of unique cards
		if hands[i].unique > hands[j].unique {
			return true
		}
		// hand is "less" when is has the same number of unique cards but lower number of repeat cards
		if hands[i].unique == hands[j].unique && hands[i].repeat < hands[j].repeat {
			return true
		}
		// hands are the same type, must compare cards in order
		if hands[i].unique == hands[j].unique && hands[i].repeat == hands[j].repeat {
			for k := range hands[i].Hand {
				ih := getOrder(string(hands[i].Hand[k]))
				jh := getOrder(string(hands[j].Hand[k]))
				if ih == jh {
					continue
				}
				if ih > jh {
					return true
				} else {
					return false
				}
			}
		}
		return false
	})
	log.Printf("Ranked hands: %v", hands)

	var winnings int
	for rank, hand := range hands {
		winnings += (rank + 1) * hand.Bid
	}
	log.Printf("Total Winnings: %d", winnings)
}

// Number of unique values:
// 1 = five of a kind
// 2 = four of a kind OR full house
// 3 = three of a kind OR two pair
// 4 = one pair
// 5 = high card

// Repeat value counter:
// 1 = high card
// 2 = pair
// 3 = three of a kind
// 4 = four of a kind
// 5 = five of a kind

// Finding hand type: (unique, repeat)
// five of a kind = (1, 5)
// four of a kind = (2, 4)
// full house = (2, 3)
// three of a kind = (3, 3)
// two pairs = (3, 2)
// one pair = (4, 2)
// high card = (5, 1)

type Hand struct {
	Hand string
	Bid, unique, repeat int
}

func newHand(hand string, bid int) Hand {
	h := Hand{Hand: hand, Bid: bid}
	m := make(map[string]int)
	wildCount := 0
	for i := range hand {
		char := string(hand[i])
		if char == "J" {
			wildCount++
		} else {
			m[string(hand[i])]++
		}
		
	}
	h.unique = len(m)

	// Fix issue when hand is all Jokers
	if h.unique == 0 {
		h.unique = 1
	}

	for _, val := range m {
		if val > h.repeat {
			h.repeat = val
		}
	}
	h.repeat += wildCount
	return h
}

func getOrder(s string) int {
	for p, v := range []string {"A", "K", "Q", "T", "9", "8", "7", "6", "5", "4", "3", "2", "J"} {
		if v == s {
			return p
		}
	}
	return -1
}