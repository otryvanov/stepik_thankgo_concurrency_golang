package main

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

// counter stores the number of digits in each word.
// each key is a word and value is the number of digits in the word.
type counter map[string]int

// countDigitsInWords counts digits in pharse words
func countDigitsInWords(phrase string) counter {
	words := strings.Fields(phrase)
	syncStats := sync.Map{}

	var wg sync.WaitGroup

	// начало решения

	// Посчитайте количество цифр в словах,
	// используя отдельную горутину для каждого слова.

	// Чтобы записать результаты подсчета,
	// используйте syncStats.Store(word, count)

	// В результате syncStats должна содержать слова
	// и количество цифр в каждом.
	wg.Add(len(words))
	for _, world := range words {
		go func(w string) {
			defer wg.Done()
			count := 0
			for _, v := range w {
				if unicode.IsDigit(v) {
					count++
				}
			}
			syncStats.Store(w, count)

		}(world)
	}
	wg.Wait()
	// конец решения

	return asStats(&syncStats)
}

// countDigits returns the number of digits in a string
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

// asStats converts stats from sync.Map to ordinary map
func asStats(m *sync.Map) counter {
	stats := counter{}
	m.Range(func(word, count any) bool {
		stats[word.(string)] = count.(int)
		return true
	})
	return stats
}

// printStats prints words and their digit counts
func printStats(stats counter) {
	for word, count := range stats {
		fmt.Printf("%s: %d\n", word, count)
	}
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	counts := countDigitsInWords(phrase)
	printStats(counts)
}
