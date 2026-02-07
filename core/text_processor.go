package core

import (
	"strings"
	"unicode"
)

// Removes extra spaces from the list of lines from the txt.
// Sadly it does not remove spaces before or after punctuation marks or quotes, because it will be too big to the eye.
// This filter will be handled by a helper function.
func ExtraSpaces(lines []string) []string {
	result := make([]string, 0, len(lines))
	for i := range lines {
		spacedLine := SeparatePunctuation(lines[i])
		words := strings.Fields(spacedLine)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}
		words = ProcessTags(words)
		filtered := Punctuation(words)
		result = append(result, filtered)
	}
	return result
}

// Helper function for edge cases like "Yes,I am".
func SeparatePunctuation(line string) string {
	var result strings.Builder
	for _, char := range line {
		if strings.ContainsRune(",.!?:;\"'", char) {
			result.WriteRune(' ')
			result.WriteRune(char)
			result.WriteRune(' ')
		} else if char == '(' {
			result.WriteRune(' ')
			result.WriteRune(char)
		} else if char == ')' {
			result.WriteRune(char)
			result.WriteRune(' ')
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// Helper function to help space filtering.
func Punctuation(words []string) string {
	var result strings.Builder //Using strings.Builder for better string manipulation and speed efficiency.
	inQuote := false

	for i, word := range words {
		if i > 0 {
			wantSpace := true

			// For edge cases like like ". .." etc.
			firstChar := word[0]

			// Checking if we are at a punctuation using the first rune to handle the edge cases mention above.
			if strings.ContainsRune(",.!?:;", rune(firstChar)) {
				wantSpace = false
			}
			// Checking if are entering a quote.
			if i > 0 && (words[i-1] == "'" || words[i-1] == "\"") && inQuote {
				wantSpace = false
			}

			// Checking if we are closing a quote.
			if (word == "'" || word == "\"") && inQuote {
				wantSpace = false
			}
			// If the current word is the apostrophe in "don't", we don't want a space before it.
			if IsContraction(words, i) {
				wantSpace = false
			}
			// If previous word was the apostrophe in "don't", we don't want a space after it.
			if IsContraction(words, i-1) {
				wantSpace = false
			} // Abbreviations check
			if i > 1 && words[i-1] == "." && len(words[i-2]) == 1 {
				if len(word) == 1 {
					// Next char must be a dot
					nextIsDot := i+1 < len(words) && words[i+1] == "."

					// Case must match (Prevent "A. b." lists or "U.S. a")
					prevLetter := []rune(words[i-2])[0]
					currLetter := []rune(word)[0]
					caseMatches := unicode.IsUpper(prevLetter) == unicode.IsUpper(currLetter)

					if nextIsDot && caseMatches {
						wantSpace = false
					}
				}
			}
			// If the character is not a punctuation, we want a space.
			if wantSpace {
				result.WriteString(" ")
			}
			// If there is a word after this one we can check if the current one is an article.
			if i+1 < len(words) {
				word = IndefiniteArticle(word, words[i+1])
			}
		}

		result.WriteString(word)

		if word == "'" || word == "\"" {
			inQuote = !inQuote
		}
	}
	return result.String()
}

// This function will check if the single quote is a quote or a contraction.
func IsContraction(words []string, i int) bool {
	//Needs a word before and a word after.
	if i == 0 || i >= len(words)-1 {
		return false
	}
	// Only single quote can be used as contraction.
	if words[i] != "'" {
		return false
	}

	prevWord := words[i-1]
	nextWord := words[i+1]

	// Checking if there are any punctuations in previous word or in the next one to avoid edge cases like "Hello , ' world".
	if strings.ContainsAny(prevWord, ",.!?:;\"'") || strings.ContainsAny(nextWord, ",.!?:;\"'") {
		return false
	}
	// Checking if the next word matches any words that are used for contraction.
	switch strings.ToLower(nextWord) {
	case "t", "s", "d", "m", "ll", "re", "ve":
		return true
	default:
		return false
	}
}

// This function searches and changes incorrect indefinite articles.
func IndefiniteArticle(word, nextword string) string {
	lowerword := strings.ToLower(word)
	if lowerword != "a" && lowerword != "an" {
		return word
	}
	// If there is no next word then there is no need to change it.
	if len(nextword) == 0 {
		return word
	}

	// Calling a helper function to check if the next word has a vowel sound.
	isVowel := HasVowelSound(nextword)

	if lowerword == "a" && isVowel {
		if word == "A" {
			return "An"
		}
		return "an"
	}
	if lowerword == "an" && !isVowel {
		if word == "An" {
			return "A"
		}
		return "a"
	}
	return word
}

// Helper function for indefinite article, looks for some edge cases too.
func HasVowelSound(word string) bool {
	lowerword := strings.ToLower(word)
	firstChar := lowerword[0]

	switch firstChar {
	case 'a', 'i':
		return true
	case 'e':
		if strings.HasPrefix(lowerword, "eu") {
			return false
		}
		return true
	case 'o':
		if strings.HasPrefix(lowerword, "one") || strings.HasPrefix(lowerword, "once") {
			return false
		}
		return true
	case 'u':
		consonantSoundU := []string{"uni", "use", "usa", "usu", "ute"}
		vowelSoundUni := []string{"unid", "unim", "unin", "unip", "uniq"}

		isConsonant := false
		for _, c := range consonantSoundU {
			if strings.HasPrefix(lowerword, c) {
				isConsonant = true
				break
			}
		}

		if isConsonant {
			// Check if it's actually one of the "unidentified" cases.
			for _, v := range vowelSoundUni {
				if strings.HasPrefix(lowerword, v) {
					return true
				}
			}
			return false // It is a consonant sound (Unique, University).
		}
		return true // Default 'u' is vowel (Umbrella, Ugly).
	case 'h':
		if strings.HasPrefix(lowerword, "honest") ||
			strings.HasPrefix(lowerword, "hour") ||
			strings.HasPrefix(lowerword, "heir") ||
			strings.HasPrefix(lowerword, "honor") {
			return true
		}
		return false
	default:
		return false
	}
}
