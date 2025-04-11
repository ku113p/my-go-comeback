package exercise3

import "strings"

func WordCount(s string) map[string]int {
	words_count := map[string]int{}

	s_words := strings.Split(s, " ")

	for _, w := range s_words {
		if w == "" {
			continue
		}
		_, ok := words_count[w]
		if ok {
			words_count[w]++
		} else {
			words_count[w] = 1
		}
	}

	return words_count
}
