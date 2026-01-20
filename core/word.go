package core

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// We cannot use a simple string for a key
// as spanish word pan (for bread) and czech word pan (for sir) would map a same key, which is undesirable
type WordId struct {
	Word string
	// If the word is in the new language I am learning
	NewLang bool
}

func (w WordId) String() string {
	if w.NewLang {
		return w.Word + "_from"
	} else {
		return w.Word + "_to"
	}
}

type Word struct {
	Id   WordId
	From string
	To   []string
}

func (i Word) Test(input string) bool {
	for _, toVal := range i.To {
		if compare(input, toVal) {
			return true
		}
	}
	return false
}

func compare(actual, expected string) bool {
	e, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(expected)))
	a, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(actual)))
	return e == a
}

var transformer = transform.Chain(
	// Decompose
	norm.NFD,

	runes.Remove(runes.Predicate(func(r rune) bool {
		// Is nonspacing mark? Which means that the combined character doesn't take more horizontal space
		return unicode.Is(unicode.Mn, r)
	})),
	// Compose
	norm.NFC,
)
