package main

import "log"

type Module interface {
	Pulse(from Module, queue *MessageQueue, high bool, iterations int)
	IsOriginalState() bool
	AddInput(input Module)
	AddOutput(output Module)
	IsRX() bool
}

// Flip-flop
type FlipFlop struct {
	State bool
	Destinations []Module
}

func (f *FlipFlop) Pulse(_ Module, queue *MessageQueue, high bool, _ int) {
	// on low pulse
	if !high {
		f.State = !f.State
		for _, dest := range f.Destinations {
			queue.add(Message{from: f, to: dest, high: f.State})
		}
	}
}

func (f *FlipFlop) IsOriginalState() bool {
	return !f.State
}

func (f *FlipFlop) AddInput(_ Module) {
	// NOOP
}

func (f *FlipFlop) AddOutput(output Module) {
	f.Destinations = append(f.Destinations, output)
}

func (f *FlipFlop) IsRX() bool {
	return false
}

// Conjunction
type Conjunction struct {
	InputMemory map[Module]bool
	Destinations []Module
}

func (c *Conjunction) Pulse(from Module, queue *MessageQueue, high bool, iterations int) {
	c.InputMemory[from] = high
	pulse := false
	for _, high := range c.InputMemory {
		if !high {
			pulse = true
		}
	}

	for _, dest := range c.Destinations {
		// Tracking cycles when output is to rx
		if dest.IsRX() {
			queue.rxPrecursorCount = len(c.InputMemory)
			for module, high := range c.InputMemory {
				if high {
					log.Println("RX Precursor Cycle:", iterations)
					queue.trackRxCycle(module, iterations)
				}
			}
		}
		queue.add(Message{from: c, to: dest, high: pulse})
	}
}

func (c *Conjunction) IsOriginalState() bool {
	r := true
	for _, high := range c.InputMemory {
		if high {
			r = false
			break
		}
	}
	return r
}

func (c *Conjunction) AddInput(input Module) {
	c.InputMemory[input] = false
}

func (c *Conjunction) AddOutput(output Module) {
	c.Destinations = append(c.Destinations, output)
}

func (c *Conjunction) IsRX() bool {
	return false
}

// Broadcast
type Broadcast struct {
	Destinations []Module
}

func (b *Broadcast) Pulse(_ Module, queue *MessageQueue, high bool, _ int) {
	for _, dest := range b.Destinations {
		queue.add(Message{from: b, to: dest, high: high})
	}
}

func (b *Broadcast) IsOriginalState() bool {
	return true
}

func (b *Broadcast) AddInput(_ Module) {
	// NOOP
}

func (b *Broadcast) AddOutput(output Module) {
	b.Destinations = append(b.Destinations, output)
}

func (b *Broadcast) IsRX() bool {
	return false
}

// Null

type Null struct {
	label string
	state bool
}

func (n *Null) Pulse(_ Module, _ *MessageQueue, high bool, _ int) {
	n.state = high
}

func (n *Null) IsOriginalState() bool {
	return true
}

func (n *Null) AddInput(_ Module) {
	// NOOP
}

func (n *Null) AddOutput(_ Module) {
	// NOOP
}

func (n *Null) IsRX() bool {
	return n.label == "rx"
}