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

func NewStringSet(elements ...string) StringSet {
	s := make(StringSet)
	s.AddAll(elements...)
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

func (s StringSet) Clear() {
	if s == nil {
		return
	}

	for k := range s {
		delete(s, k)
	}
}

func (s StringSet) ReplaceAll(elements []string) {
	if s == nil {
		s = make(StringSet)
	}

	s.Clear()
	s.AddAll(elements...)
}

func Contains(stringSlice []string, element string) bool {
	for _, s := range stringSlice {
		if s == element {
			return true
		}
	}
	return false
}
