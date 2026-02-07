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

				lastIdx := len(result) - 1
				//Strip punctuation from the value before parsing (e.g., "1E," -> "1E").
				cleanVal := strings.TrimFunc(result[lastIdx], func(r rune) bool {
					return !unicode.IsDigit(r) && !unicode.IsLetter(r)
				})

				if val, err := strconv.ParseInt(cleanVal, base, 64); err == nil {
					result[lastIdx] = fmt.Sprint(val)
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

		//The tag is strictly split like ["(up", ",", "2)"]
		//This happens because SeparatePunctuation ran first.
		if word != "(up)" && word != "(low)" && word != "(cap)" && i+2 < len(words) {
			nextWord := words[i+1]
			numberWord := words[i+2]

			if nextWord == "," {
				// We found the comma, now try to parse the number.
				cleanNum := strings.TrimFunc(numberWord, func(r rune) bool {
					return !unicode.IsDigit(r)
				})

				if val, err := strconv.Atoi(cleanNum); err == nil {
					count = val
					i += 2 //Skip the comma and the number.
				}
			}
		}

		if strings.HasSuffix(word, ",") && i+1 < len(words) {
			nextWord := words[i+1]
			cleanNum := strings.TrimFunc(nextWord, func(r rune) bool {
				return !unicode.IsDigit(r)
			})
			if val, err := strconv.Atoi(cleanNum); err == nil {
				count = val
				i += 1 //Skip the number used as an argument.
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
