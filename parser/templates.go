package parser

import "strconv"

// Add adds up the given summands. It is used in templates
func Add(summands ...interface{}) int {
	var sum int
	for _, s := range summands {
		if i, ok := s.(int); ok {
			sum += i
		} else if s, ok := s.(string); ok {
			i, _ := strconv.Atoi(s)
			sum += i
		}
	}
	return sum
}
