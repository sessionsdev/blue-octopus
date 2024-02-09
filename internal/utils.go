package internal

type StringSet map[string]struct{}
type void struct{}

var member void

func GetStringsFromMap(m map[string]struct{}) []string {
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

func MakeSet(elements ...string) StringSet {
	s := make(StringSet)
	for _, e := range elements {
		s[e] = member
	}
	return s
}

func EmptyStringSet() StringSet {
	return make(StringSet)
}

func (s StringSet) ToSlice() []string {
	return GetStringsFromMap(s)
}

func (s StringSet) AddAll(elements ...string) {
	if s == nil {
		s = make(StringSet)
	}

	for _, e := range elements {
		s[e] = member
	}
}

func (s StringSet) RemoveAll(elements ...string) {
	if s == nil {
		return
	}

	for _, e := range elements {
		delete(s, e)
	}
}

func (s StringSet) Contains(element string) bool {
	_, ok := s[element]
	return ok
}
