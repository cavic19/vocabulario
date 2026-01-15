package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
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

const SnapshotInterval time.Duration = 10 * time.Second

func main() {
	var dataDir string
	var lessonName string

	flag.StringVar(&dataDir, "dir", "data", "Directory containing lesson CSV files")
	flag.StringVar(&lessonName, "lesson", "", "Specific lesson to load (optional, loads all if not specified)")
	flag.Parse()

	lessons, err := loadLessons(dataDir, lessonName)
	if err != nil {
		log.Fatalf("Error loading lessons: %v", err)
	}

	if len(lessons) == 0 {
		fmt.Println("No lessons found.")
		return
	}

	var stats map[VocabularyPair]SuccessStats = make(map[VocabularyPair]SuccessStats)
	loadedStats := LoadStats(dataDir)
	for _, lesson := range lessons {
		for _, pair := range lesson.Pairs {
			pair1 := VocabularyPair{pair.From, pair.To}
			pair2 := VocabularyPair{pair.To, pair.From}
			// Remember, when there is no value for givne pair in the loaded map, it will use a default value which is basically all zeros
			stats[pair1] = loadedStats[pair1]
			stats[pair2] = loadedStats[pair2]
		}
	}
	reader := bufio.NewReader(os.Stdin)

	lastSave := time.Now()

	for {
		if time.Since(lastSave) >= SnapshotInterval {
			SaveStats(stats, dataDir)
			lastSave = time.Now()
		}

		word := NextWord(stats)
		fmt.Printf("%v: ", word.From)
		input, _ := reader.ReadString('\n')
		oldStats := stats[word]
		if Compare(word.To, input) {
			stats[word] = SuccessStats{
				oldStats.Success + 1,
				oldStats.Failure,
			}
		} else {
			fmt.Printf("Wrong! It should be %v\n", word.To)
			stats[word] = SuccessStats{
				oldStats.Success,
				oldStats.Failure + 1,
			}
		}
	}
}

func Compare(expected, actual string) bool {
	e, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(expected)))
	a, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(actual)))
	return e == a
}
