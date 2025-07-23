// Package random is a test package
package random

import (
	"testing"
	"time"
)


func TestNewRandomStringOnTheLength(t *testing.T) {
//* Arrange
	lengths := []int{1, 2, 4, 7, 10, 70}
	
//* Act 
	for _, length := range lengths {
		lenRandString := len(NewRandomString(length))
//* Assert
		if length != lenRandString {
			t.Errorf("Expected %d length, got %d", length, lenRandString)
		}
	}
}

func TestNewRandomStringOnTheRange(t *testing.T) {
//Arrange
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
//Act
		resFunc := NewRandomString(4)
//Assert
		if validateString(resFunc, letters) {
			t.Logf("Great, string %s is valid", resFunc)
		} else {
			t.Error("Expected that string is valid")
		}
	}
}

func TestNewRandomStringOnTheNotEqual(t *testing.T) {
	lengths := []int{1, 2, 4, 7, 10, 70}
	for _, length := range lengths {
		str1 := NewRandomString(length)
		str2 := NewRandomString(length)
		if str1 == str2 {
			t.Errorf("Duplicate string found for length %d: %s", length, str1)
		}
	}

}

func validateString(s string, runes []rune) bool { //гениальная функция
	for _, char := range s {
		found := false
		for _, r := range runes {
			if char == r {
				found = true
				break
			} 
		}
		if !found {
			return false
		}
	}
	return true
}
