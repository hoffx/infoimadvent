package parser

import "strconv"

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
