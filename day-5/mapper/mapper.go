package mapper

type Mapper struct {
	DestStart int
	SrcStart int
	Length int
}

type SeedRange struct {
	Start int
	Length int
}

func (m Mapper) inSourceRange(check int) bool {
	if check >= m.SrcStart && check < m.SrcStart + m.Length {
		return true
	} else {
		return false
	}
}

func (m Mapper) getDest(src int) (int, bool) {
	if !m.inSourceRange(src) {
		return 0, false
	}

	diff := src - m.SrcStart
	return m.DestStart + diff, true
}

func Resolve(mappers []Mapper, src int) int {
	for _, mapper := range mappers {
		dest, matched := mapper.getDest(src)
		if matched {
			return dest
		}
	}
	return src
}