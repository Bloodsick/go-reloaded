package core

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func ProcessTags(words []string) []string {
	result := make([]string, 0, len(words))

	for i := 0; i < len(words); i++ {
		word := words[i]

		//Identify the transformation we want to use.
		var transformFunc func(string) string

		//Handle Hex/Bin cases since they are easy to find and transform.
		if word == "(hex)" || word == "(bin)" {
			if len(result) > 0 {
				base := 10
				if word == "(hex)" {
					base = 16
				} else {
					base = 2
				}

				// We look backwards to find a valid word to convert.
				targetIdx := -1
				for j := len(result) - 1; j >= 0; j-- {
					if !IsPunctuation(result[j]) {
						targetIdx = j
						break
					}
				}

				// If we found a valid word to convert.
				if targetIdx != -1 {
					// Strip punctuation from the value before parsing.
					cleanVal := strings.TrimFunc(result[targetIdx], func(r rune) bool {
						return !unicode.IsDigit(r) && !unicode.IsLetter(r)
					})

					if val, err := strconv.ParseInt(cleanVal, base, 64); err == nil {
						result[targetIdx] = fmt.Sprint(val)
					}
				}
			}
			continue
		}

		//Tags.
		switch {
		case strings.HasPrefix(word, "(up"):
			transformFunc = strings.ToUpper
		case strings.HasPrefix(word, "(low"):
			transformFunc = strings.ToLower
		case strings.HasPrefix(word, "(cap"):
			transformFunc = Capitalize
		default:
			result = append(result, word)
			continue
		}

		count := 1
		argIdx := i + 1

		// Skip the comma if it's separated.
		if argIdx < len(words) && words[argIdx] == "," {
			argIdx++
		}

		// Parse the number argument.
		if argIdx < len(words) {
			argWord := words[argIdx]
			cleanNum := strings.TrimFunc(argWord, func(r rune) bool {
				return !unicode.IsDigit(r)
			})

			if val, err := strconv.Atoi(cleanNum); err == nil && val > 0 {
				count = val
				i = argIdx // Advance main loop to the number.

				// Consume trailing ')' if it's separated.
				if i+1 < len(words) && words[i+1] == ")" {
					i++
				}
			}
		}

		//Apply the transformation backwards.
		for j := 0; j < count; j++ {
			idx := len(result) - 1 - j
			if idx < 0 {
				break
			}
			result[idx] = transformFunc(result[idx])
		}
	}

	return result
}

// Helper function to check if the string only has punctuation characters.
func IsPunctuation(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !strings.ContainsRune(",.!?:;\"'()", r) {
			return false
		}
	}
	return true
}

// Helper for Capitalize
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	for i := 1; i < len(r); i++ {
		r[i] = unicode.ToLower(r[i])
	}
	return string(r)
}
