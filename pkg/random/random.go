package random

import (
	"math/rand"
	"time"
)

const (
	smallLetter = "abcdefghijklmnopqrstuvwxyz"
	capLetter   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitLetter = "0123456789"
)

func String(size int, small, cap, digit bool) string {
	dict := ""
	if small {
		dict = dict + smallLetter
	}
	if cap {
		dict = dict + capLetter
	}
	if digit {
		dict = dict + digitLetter
	}
	letters := []rune(dict)

	rand.Seed(time.Now().UnixNano())
	s := make([]rune, size)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
