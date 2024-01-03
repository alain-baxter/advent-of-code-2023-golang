package node

import "regexp"

type Node struct {
	Label string
	Left string
	Right string
}

const parser string = "^([0-9A-Z]{3})\\s=\\s\\(([0-9A-Z]{3}),\\s([0-9A-Z]{3})"

func NewNode(s string) Node {
	var n Node
	r := regexp.MustCompile(parser)
	match := r.FindStringSubmatch(s)
	n.Label = match[1]
	n.Left = match[2]
	n.Right = match[3]
	return n
}