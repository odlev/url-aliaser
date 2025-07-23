// Package random is a nice package
package random

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewRandomString(length int) string {

	symbols := make([]rune, length)

	for i := range symbols {
		symbols[i] = letters[seed.Intn(len(letters))]
	}

	return string(symbols)
}

