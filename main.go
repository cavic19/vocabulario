package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type VocabularyPair struct {
	From string
	To   string
}

type Lesson struct {
	Name  string
	Pairs []VocabularyPair
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

func main() {
	var dataDir string
	var lessonName string

	flag.StringVar(&dataDir, "dir", "data", "Directory containing lesson CSV files")
	flag.StringVar(&lessonName, "lesson", "", "Specific lesson to load (optional, loads all if not specified)")
	flag.Parse()

	lessons, err := LoadLessons(dataDir, lessonName)
	if err != nil {
		log.Fatalf("Error loading lessons: %v", err)
	}

	if len(lessons) == 0 {
		fmt.Println("No lessons found.")
		return
	}

	var stats *WordStats = WordStatsFromFile(dataDir)
	for _, lesson := range lessons {
		for _, pair := range lesson.Pairs {
			pair1 := VocabularyPair{pair.From, pair.To}
			pair2 := VocabularyPair{pair.To, pair.From}
			if _, ok := stats.counts[pair1]; !ok {
				stats.counts[pair1] = WordRecord{}
			}
			if _, ok := stats.counts[pair2]; !ok {
				stats.counts[pair2] = WordRecord{}
			}
		}
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		SaveStats(stats, dataDir)
		word := stats.NextWord()
		fmt.Printf("%v (%.2f): ", word.From, stats.GetStats(word).SuccessRate())
		input, _ := reader.ReadString('\n')
		if Compare(word.To, input) {
			stats.RecordSuccess(word)
		} else {
			fmt.Printf("Wrong! It should be %v\n", word.To)
			stats.RecordFailure(word)
		}
	}
}

func Compare(expected, actual string) bool {
	e, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(expected)))
	a, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(actual)))
	return e == a
}
