package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
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
	var matchingSum uint64 = 0
	folds := 4
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Fields(text)
		pattern := parsePattern(fields[0], folds)
		groups := parseGroups(fields[1], folds)
		log.Printf("Pattern: %s Groupings: %d", pattern, groups)
		//count := countRecursionGroupMatching(pattern, &groups)
		//count := countRecursionRegex(pattern, regexp.MustCompile(createRegex(groups)))
		count := countIteration(pattern, groups)
		log.Printf("Number of matching options: %d", count)
		matchingSum += uint64(count)
	}
	log.Printf("Sum of matches: %d", matchingSum)
}

func parsePattern(s string, folds int) string {
	r := s
	for i := 0; i < folds; i++ {
		r = r + "?" + s
	}
	return r
}

func parseGroups(s string, folds int) []int {
	var r []int
	fields := strings.Split(s, ",")
	for _, field := range fields {
		val, _ := strconv.Atoi(field)
		r = append(r, val)
	}
	iter := slices.Clone(r)
	for i := 0; i < folds; i++ {
		r = append(r, iter...)
	}
	return r
}

/*
  This method is the first algorithm at counting the options I created, and was used for part 1.
	It is recursive, going down both branches when we find a "?" and computes the groups and compares 
	to the target groups.

	When starting part 2 I added several improvements like early detection that a branch was not 
	matching to prune recursion branches. I tried to multi-thread the two branches, but this causes
	OOM issues with the thread overhead.
*/
func countRecursionGroupMatching(pattern string, target *[]int) int {
	count := 0

	// Check if we have a group mismatch as we recurse, leave branches early
	groups := computeGroups(pattern)
	for i, group := range groups {
		// If we have too many groups leave early
		if i >= len(*target) {
			return 0
		}
		// If the current group is bigger than the corresponding target group, leave early
		if group > (*target)[i] {
			return 0
		}
		// if we aren't the last group and the group is not exactly the corresponding target group, leave early
		if i < len(groups) - 1 && group != (*target)[i] {
			return 0
		}
	}

	unknown := strings.Index(pattern, "?")
	if unknown == -1 {
		// No more unknown entries, check for exact group match (all but last entry has been verified)
		if len(groups) == len(*target) && groups[len(groups) - 1] == (*target)[len(*target) - 1] {
			return 1
		}
		return 0
	} else {
		count += countRecursionGroupMatching(strings.Replace(pattern, "?", ".", 1), target)
		count += countRecursionGroupMatching(strings.Replace(pattern, "?", "#", 1), target)
	}
	return count
}

func computeGroups(pattern string) []int {
	var groups []int
	buffer := 0
	for i := range pattern {
		val := string(pattern[i])
		if buffer != 0 && val == "." {
			groups = append(groups, buffer)
			buffer = 0
		}
		if val == "#" {
			buffer++
		}
		if val == "?" {
			break
		}
	}

	if buffer != 0 {
		groups = append(groups, buffer)
	}

	return groups
}

/*
  This method is similar to countRecursionGroupMatching in that it is a recursive
	algorithm, however I was hoping that using a dynamic regex created by createRegex
	to simplify and speed up the recursion. 
	
	It was a bit more performant but nowhere near what was needed for part 2 to run 
	successfully within a weekend (got to 90/1000)
*/
func countRecursionRegex(pattern string, matcher *regexp.Regexp) uint64 {
	var count uint64 = 0

	if matcher.MatchString(pattern) {
		if !strings.Contains(pattern, "?") {
			//log.Printf("Matched: %s", pattern)
			return 1
		}
		count += countRecursionRegex(strings.Replace(pattern, "?", ".", 1), matcher)
		count += countRecursionRegex(strings.Replace(pattern, "?", "#", 1), matcher)
	}

	return count
}

func createRegex(target []int) string {
	regex := "^[\\.|?]*"

	for i, group := range target {
		if i != 0 {
			regex += "[\\.|?]+"
		}
		regex += "[#|?]{" + fmt.Sprint(group) + "}"
	}

	regex += "[\\.|?]*$"

	return regex
}

/*
  This method was required to do part 2 in any reasonable amount of comupte time.

	This iterates over the pattern character by character and keeps track of the valid
	combinations in a map. We drop any invalid options along the way. This is iterative
	as opposed to recursive like the last 2 attempts and runs significantly faster.

	It does part 2 in 1 second.
*/
func countIteration(pattern string, groups []int) uint64 {
	tracker := make(map[Tracker]int)
	tracker[Tracker{0, 0}] = 0

	for i := range pattern {
		val := string(pattern[i])
		//log.Printf("Checking character: %s", val)
		trackerAdjustments := make(map[Tracker]int)

		for key, count := range tracker {

			if val == "#" {
				// It is a "#". Verify the current group can be expanded by 1
				if isValid(key.currentSize + 1, key.groupIndex, groups, false) {
				  addTracker(&trackerAdjustments, key.groupIndex, key.currentSize + 1, count)
			  }

			} else if val == "?" {
				// What if it were a "."? It ends the current group so make sure it's the right size
				if key.currentSize > 0 {
					if isValid(key.currentSize, key.groupIndex, groups, true) {
					  addTracker(&trackerAdjustments, key.groupIndex + 1, 0, count)
					}
				} else {
					// Keep the entry if we don't have a current group counting
					addTracker(&trackerAdjustments, key.groupIndex, key.currentSize, count)
				}

				// What if it were a "#"? It increments the current group size if allowed
				if isValid(key.currentSize + 1, key.groupIndex, groups, false) {
					addTracker(&trackerAdjustments, key.groupIndex, key.currentSize + 1, count)
				}

			} else {
				// It is a ".". Verify the current group can be ended
				if key.currentSize > 0 {
					if isValid(key.currentSize, key.groupIndex, groups, true) {
				    addTracker(&trackerAdjustments, key.groupIndex + 1, 0, count)
			    }
				} else {
					addTracker(&trackerAdjustments, key.groupIndex, key.currentSize, count)
				}
			}
		}

		tracker = trackerAdjustments
		//log.Printf("Incremental tracker: %v", tracker)
	}

	//log.Printf("Final tracker: %v", tracker)
	var r uint64 = 0
	for key, count := range tracker {
		// Add the tracker that made it to the correct groupings right at the end
		if key.groupIndex == len(groups) - 1 && key.currentSize == groups[key.groupIndex] {
			r += uint64(count)
		}
		// Add the tracker that made it to the correct groupings early
		if key.groupIndex == len(groups) && key.currentSize == 0 {
			r += uint64(count)
		}
	}
	return r
}

func isValid(currentSize int, groupIndex int, groups []int, strict bool) bool {
	if groupIndex >= len(groups) {
		return false
	}

	if strict {
		return currentSize == groups[groupIndex]
	} else {
		return currentSize <= groups[groupIndex]
	}
}

func addTracker(tracker *map[Tracker]int, groupIndex int, currentGroup int, count int) {
	if count == 0 {
		(*tracker)[Tracker{groupIndex, currentGroup}]++
	} else {
		(*tracker)[Tracker{groupIndex, currentGroup}] += count
	}
}

type Tracker struct {
	groupIndex, currentSize int
}