package main

import (
	"fmt"
)

func isValid(s string) bool {

	if len(s) < 1 || len(s) > 4096 {
		return false
	}

	stack := []rune{}
	pairs := map[rune]rune{
		'>': '<',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		if char == '<' || char == '{' || char == '[' {
			stack = append(stack, char)
		} else if char == '>' || char == '}' || char == ']' {
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			if pairs[char] != top {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			return false
		}
	}

	return len(stack) == 0
}

func main() {
	var text string
	fmt.Print("input text : ")
	fmt.Scan(&text)

	fmt.Println(isValid(text))
	return
}
