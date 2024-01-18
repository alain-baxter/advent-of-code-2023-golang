package main

import "alain-baxter/aoc-2023/utils"

// FIFO queue
type MessageQueue struct {
	entries           []Message
	// For part 1
	lowCount          int
	highCount         int
	// For part 2
	rxPrecursorCount  int
	rxPrecursorCycles map[Module]int
}

type Message struct {
	from, to Module
	high     bool
}

func (q *MessageQueue) add(msg Message) {
	q.entries = append(q.entries, msg)
	if msg.high {
		q.highCount++
	} else {
		q.lowCount++
	}
}

func (q *MessageQueue) pop() Message {
	msg := q.entries[0]
	q.entries = q.entries[1:]
	return msg
}

func (q *MessageQueue) trackRxCycle(precursor Module, cycle int) {
	q.rxPrecursorCycles[precursor] = cycle
}

// Another problem where the solution is built from the specific input.
// rx is sent a message from a conjunction, so it only gets a low when
// all the input modules to the conjunction are high. They will cycle
// and only line up at the LCM of their cycles
func (q MessageQueue) checkRxPrecursorCycles() (int, bool) {
	count := 0
	product := 1
	for _, cycle := range q.rxPrecursorCycles {
		product = utils.Lcm(product, cycle)
		count++
	}

	if q.rxPrecursorCount > 0 && count == q.rxPrecursorCount {
		return product, true
	}

	return 0, false
}

func (q MessageQueue) execute(broadcast Broadcast, iterations int) int {
	for i := 0; i < iterations; i++ {
		q.add(Message{from: nil, to: &broadcast, high: false})

		for {
			msg := q.pop()
			msg.to.Pulse(msg.from, &q, msg.high, i)

			if len(q.entries) == 0 {
				break
			}
		}
	}

	return q.lowCount * q.highCount
}

func (q MessageQueue) toRX(broadcast Broadcast) int {
	for i := 1; ; i++ {
		q.add(Message{from: nil, to: &broadcast, high: false})

		for {
			msg := q.pop()
			msg.to.Pulse(msg.from, &q, msg.high, i)

			if msg.to.IsRX() {
				r, ok := q.checkRxPrecursorCycles()
				if ok {
					return r
				}
			}

			if len(q.entries) == 0 {
				break
			}
		}
	}
}
