package main

import (
	"alain-baxter/aoc-2023/day-5/mapper"
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

	var sd2sl, sl2fz, fz2wr, wr2lt, lt2tp, tp2hm, hm2ln []mapper.Mapper
	var seeds []mapper.SeedRange
	var currMapper *[]mapper.Mapper
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		switch {
		case strings.Contains(text, "seeds:"):
			seeds = getSeeds(text)
			continue
		case strings.Contains(text, "seed-to-soil"):
			currMapper = &sd2sl
			continue
		case strings.Contains(text, "soil-to-fertilizer"):
			currMapper = &sl2fz
			continue
		case strings.Contains(text, "fertilizer-to-water"):
			currMapper = &fz2wr
			continue
		case strings.Contains(text, "water-to-light"):
			currMapper = &wr2lt
			continue
		case strings.Contains(text, "light-to-temperature"):
			currMapper = &lt2tp
			continue
		case strings.Contains(text, "temperature-to-humidity"):
			currMapper = &tp2hm
			continue
		case strings.Contains(text, "humidity-to-location"):
			currMapper = &hm2ln
			continue
		case len(text) == 0:
			continue
		}

		entry := strings.Split(text, " ")
		destStart, _ := strconv.Atoi(strings.TrimSpace(entry[0]))
		srcStart, _ := strconv.Atoi(strings.TrimSpace(entry[1]))
		length, _ := strconv.Atoi(strings.TrimSpace(entry[2]))

		*currMapper = append(*currMapper, mapper.Mapper{DestStart: destStart, SrcStart: srcStart, Length: length})
	}

	var minLocation int = 2147483647
	for _, seedRange := range seeds {
		log.Printf("Seed Range: %v", seedRange)
		for i := 0; i < seedRange.Length; i++ {
			seed := seedRange.Start + i
			soil := mapper.Resolve(sd2sl, seed)
			fertilizer := mapper.Resolve(sl2fz, soil)
			water := mapper.Resolve(fz2wr, fertilizer)
			light := mapper.Resolve(wr2lt, water)
			temperature := mapper.Resolve(lt2tp, light)
			humidity := mapper.Resolve(tp2hm, temperature)
			location := mapper.Resolve(hm2ln, humidity)
			if (location < minLocation) {
				minLocation = location
			}
			//log.Printf("Seed: %d Soil: %d Fertilizer: %d Water: %d Light: %d Temperature: %d Humidity: %d Location: %d", 
			//			seed, soil, fertilizer, water, light, temperature, humidity, location)
		}
	}

	log.Printf("Lowest Location: %d", minLocation)
}

func getSeeds(text string) []mapper.SeedRange {
	var nums []mapper.SeedRange
  var currNum int = -1 
	for _, val := range strings.Split(text, " ") {
		if (!strings.Contains(val, "seeds")) {
			if (currNum == -1) {
				currNum, _ = strconv.Atoi(strings.TrimSpace(val))
				continue
			}
			length, _ := strconv.Atoi(strings.TrimSpace(val))
			nums = append(nums, mapper.SeedRange{Start: currNum, Length: length})
			currNum = -1
		}
	}
	return nums
}