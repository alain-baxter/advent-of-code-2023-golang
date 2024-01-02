package race

const startingSpeed int = 0
const speedIncrementer int = 1

func BeatRecord(time int, record int) int {
	count := 0
	for press := 1; press < time; press++ {
		speed := startingSpeed + speedIncrementer * press
		distance := speed * (time - press)
		if distance > record {
			count++
		}
	}
	return count
}