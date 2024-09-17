package main

import (
	"fmt"
	"strings"
)

type Matcher interface {
	Match(word, dictionaryWord string) bool
	GetMatchType() string
}

type ExactMatch struct {
}

func (e *ExactMatch) Match(word, dictionaryWord string) bool {
	return word == dictionaryWord
}

func (e *ExactMatch) GetMatchType() string {
	return "Exact Match"
}

type PrefixMatch struct{}

func (p *PrefixMatch) Match(word, dictionaryWord string) bool {
	return strings.HasPrefix(dictionaryWord, word)
}

func (p *PrefixMatch) GetMatchType() string {
	return "Prefix match"
}

type SuffixMatch struct{}

func (s *SuffixMatch) Match(word, dictionaryWord string) bool {
	return strings.HasSuffix(dictionaryWord, word)
}

func (s *SuffixMatch) GetMatchType() string {
	return "Suffix match"
}

type Dictionary struct {
	words   []string
	matcher Matcher
}

func (d *Dictionary) SetMatcher(matcher Matcher) {
	d.matcher = matcher
}

func NewDictionary(words []string) *Dictionary {
	return &Dictionary{words: words}
}

func (d *Dictionary) FindMatch(word string) (string, string) {
	for _, dictWord := range d.words {
		if d.matcher.Match(word, dictWord) {
			return dictWord, d.matcher.GetMatchType()
		}
	}
	return "", "No match"
}

func main() {
	// Sample dictionary
	words := []string{"book", "conspiracy", "table", "cook", "fable"}
	dictionary := NewDictionary(words)

	// Words to find matches for
	searchWords := []string{"bool", "conspiracy", "fable"}

	for _, word := range searchWords {
		// Try exact match
		dictionary.SetMatcher(&ExactMatch{})
		matchedWord, matchType := dictionary.FindMatch(word)
		if matchedWord != "" {
			fmt.Printf("Word: %s, Match: %s, Type: %s\n", word, matchedWord, matchType)
			continue
		}

		// Try prefix match
		dictionary.SetMatcher(&PrefixMatch{})
		matchedWord, matchType = dictionary.FindMatch(word)

		if matchedWord != "" {
			fmt.Printf("Word: %s, Match: %s, Type: %s\n", word, matchedWord, matchType)
			continue
		}

		// Try suffix match
		dictionary.SetMatcher(&SuffixMatch{})
		matchedWord, matchType = dictionary.FindMatch(word)
		if matchedWord != "" {
			fmt.Printf("Word: %s, Match: %s, Type: %s\n", word, matchedWord, matchType)
		} else {
			fmt.Printf("Word: %s, Match: No match found\n", word)
		}
	}
}
