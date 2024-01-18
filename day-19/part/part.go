package part

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Part struct {
	X, M, A, S int
}

type Workflow struct {
	Label string
	Rules []Rule
}

type Rule struct {
	On, Op     string
	Comparitor int
	Dest       string
}

const GREATER string = ">"
const LESS string = "<"

const EXTREME string = "x"
const MUSICAL string = "m"
const AERODYNAMIC string = "a"
const SHINY string = "s"

const ACCEPTED string = "A"
const REJECTED string = "R"

const MAX_RATING int = 4000

// in{s<1351:px,qqz}
func ParseWorkflow(text string) Workflow {
	parts := strings.Split(text, "{")
	label := parts[0]
	rules := parseRules(parts[1][:len(parts[1])-1])

	return Workflow{label, rules}
}

// s<1351:px,qqz
func parseRules(text string) []Rule {
	var rules []Rule
	r := regexp.MustCompile("^([x|m|a|s]{1})([<|>]{1})([0-9]+):(.*)$")
	for _, rule := range strings.Split(text, ",") {
		// A compare style rule
		if r.MatchString(rule) {
			groups := r.FindStringSubmatch(rule)
			comparitor, _ := strconv.Atoi(groups[3])
			rules = append(rules, Rule{groups[1], groups[2], comparitor, groups[4]})

			// A testless rule
		} else {
			rules = append(rules, Rule{Dest: rule})
		}
	}
	return rules
}

// {x=787,m=2655,a=1222,s=2876}
func ParsePart(text string) Part {
	parts := strings.Split(text[1:len(text)-1], ",")
	x, _ := strconv.Atoi(parts[0][2:])
	m, _ := strconv.Atoi(parts[1][2:])
	a, _ := strconv.Atoi(parts[2][2:])
	s, _ := strconv.Atoi(parts[3][2:])
	return Part{x, m, a, s}
}

func (r Rule) HasCondition() bool {
	return len(r.On) != 0 && len(r.Op) != 0
}

func (r Rule) Accepts() bool {
	return r.Dest == "A"
}

func (r Rule) Rejects() bool {
	return r.Dest == "R"
}

func (r Rule) Test(p Part) (string, bool) {
	if !r.HasCondition() {
		return r.Dest, true
	}

	var test int
	switch r.On {
	case EXTREME:
		test = p.X
	case MUSICAL:
		test = p.M
	case AERODYNAMIC:
		test = p.A
	case SHINY:
		test = p.S
	}

	// Greater than rule
	if r.Op == GREATER && test > r.Comparitor {
		return r.Dest, true

		// Less than rule
	} else if r.Op == LESS && test < r.Comparitor {
		return r.Dest, true

		// Doesn't match
	} else {
		return "", false
	}
}

func (w Workflow) GetNext(p Part, workflows map[string]Workflow) Workflow {
	for _, rule := range w.Rules {
		dest, match := rule.Test(p)
		if match {
			return workflows[dest]
		}
	}

	// Should not get here, all workflows end in a testless rule
	log.Fatalf("Workflow %v could not match next from %v", w, p)
	return Workflow{}
}

func CheckPart(p Part, workflows map[string]Workflow) int {
	current := workflows["in"]
	r := 0
	for {
		current = current.GetNext(p, workflows)

		if current.Label == "A" {
			r = p.getRating()
			break
		}

		if current.Label == "R" {
			break
		}
	}

	return r
}

func (p Part) getRating() int {
	return p.X + p.M + p.A + p.S
}

func BrutePathsToAccepted(workflows map[string]Workflow) uint64 {
	var accepted uint64 = 0
	for part := range generateParts() {
		if CheckPart(part, workflows) > 0 {
			accepted++
		}
	}
	return accepted
}

func generateParts() <-chan Part {
	c := make(chan Part)

	go func(c chan Part) {
		defer close(c)
		for x := 1; x <= int(MAX_RATING); x++ {
			for m := 1; m <= int(MAX_RATING); m++ {
				for a := 1; a <= int(MAX_RATING); a++ {
					for s := 1; s <= int(MAX_RATING); s++ {
						c <- Part{x, m, a, s}
					}
				}
			}
		}
	}(c)

	return c
}

type Tracker struct {
	x, m, a, s Range
}

type Range struct {
	start, end int
}

func PathsToAccepted(workflows map[string]Workflow) uint64 {
	// Always start from "in" with all possible combinations
	current := workflows["in"]
	tracker := Tracker{Range{1, MAX_RATING}, Range{1, MAX_RATING}, Range{1, MAX_RATING}, Range{1, MAX_RATING}}

	return findCombinations(current, workflows, tracker)
}

func findCombinations(workflow Workflow, workflows map[string]Workflow, ranges Tracker) uint64 {
	result := uint64(0)
	forward := ranges
	for _, rule := range workflow.Rules {
		match := forward

		if rule.HasCondition() {
			forward, match = splitRange(rule, forward)
		}

		if rule.Accepts() {
			result += getCombinations(match)
		} else if rule.Rejects() {
			result += 0
		} else {
			result += findCombinations(workflows[rule.Dest], workflows, match)
		}
	}

	return result
}

func getCombinations(ranges Tracker) uint64 {
	log.Println("Accepting:", ranges)
	return uint64(ranges.x.end-ranges.x.start+1) *
		uint64(ranges.m.end-ranges.m.start+1) *
		uint64(ranges.a.end-ranges.a.start+1) *
		uint64(ranges.s.end-ranges.s.start+1)
}

func splitRange(rule Rule, ranges Tracker) (Tracker, Tracker) {
	split := ranges

	switch rule.On {
	case EXTREME:
		split.x, ranges.x = splitAtribute(rule, ranges.x)
	case MUSICAL:
		split.m, ranges.m = splitAtribute(rule, ranges.m)
	case AERODYNAMIC:
		split.a, ranges.a = splitAtribute(rule, ranges.a)
	case SHINY:
		split.s, ranges.s = splitAtribute(rule, ranges.s)
	}

	return split, ranges
}

func splitAtribute(rule Rule, attributeRange Range) (Range, Range) {
	start := attributeRange.start
	end := attributeRange.end

	if rule.Op == GREATER {
		return Range{start, rule.Comparitor}, Range{rule.Comparitor + 1, end}
	} else {
		return Range{rule.Comparitor, end}, Range{start, rule.Comparitor - 1}
	}
}
