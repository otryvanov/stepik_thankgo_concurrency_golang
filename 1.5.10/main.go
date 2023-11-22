package main

import (
	"fmt"
	"math/rand"
)

// начало решения
type chanelString struct {
	first string
	last  string
}

// генерит случайные слова из 5 букв
// с помощью randomWord(5)
func generate(cancel <-chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-cancel:
				return
			case out <- randomWord(5):
			}
		}
	}()
	return out
}

// выбирает слова, в которых не повторяются буквы,
// abcde - подходит
// abcda - не подходит
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for in != nil {
			select {
			case <-cancel:
				return
			case v, ok := <-in:
				if !ok {
					in = nil
					return
				}

				if isIsogram(v) {
					select {
					case <-cancel:
						return
					case out <- v:
					}
				}

			}
		}
	}()
	return out
}

func isIsogram(word string) bool {
	for i := 0; i < len(word); i++ {
		for j := i + 1; j < len(word); j++ {
			if word[i] == word[j] {
				return false
			}
		}
	}

	return true
}

// переворачивает слова
// abcde -> edcba
func reverse(cancel <-chan struct{}, in <-chan string) <-chan chanelString {
	out := make(chan chanelString)
	go func() {
		defer close(out)
		for in != nil {
			select {
			case <-cancel:
				return
			case word, ok := <-in:
				if !ok {
					in = nil
					return
				}
				buf := chanelString{
					first: word,
					last:  reverseWord(word),
				}
				select {
				case <-cancel:
					return
				case out <- buf:
				}

			}
		}
	}()
	return out
}

func reverseWord(word string) string {
	runes := []rune(word)
	length := len(runes)
	for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan chanelString) <-chan chanelString {
	out := make(chan chanelString)
	go func() {
		defer close(out)
		for c1 != nil || c2 != nil {
			select {
			case <-cancel:
				return
			case v, ok := <-c1:
				if !ok {
					c1 = nil
					continue
				}
				select {
				case <-cancel:
					return
				case out <- v:
				}
			case v, ok := <-c2:
				if !ok {
					c2 = nil
					continue
				}
				select {
				case <-cancel:
					return
				case out <- v:
				}

			}
		}
	}()
	return out
}

// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan chanelString, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-cancel:
			return
		case v, ok := <-in:
			if !ok {
				in = nil
				return
			}
			fmt.Printf("%s -> %s\n", v.first, v.last)
		}
	}
}

// конец решения

// генерит случайное слово из n букв
func randomWord(n int) string {
	const letters = "aeiourtnsl"
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 10)
}
