package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
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
	var dataDir, lessonName string
	var showStats bool

	flag.StringVar(&dataDir, "dir", "data", "Directory containing lesson CSV files")
	flag.StringVar(&lessonName, "lesson", "", "Specific lesson to load (optional, loads all if not specified)")
	flag.BoolVar(&showStats, "stats", false, "Display vocabulary statistics sorted by success rate (highest to lowest)")
	flag.Parse()

	lessons, err := LoadLessons(dataDir, lessonName)
	if err != nil {
		log.Fatalf("Error loading lessons: %v", err)
	}

	if len(lessons) == 0 {
		fmt.Println("No lessons found.")
		return
	}

	if showStats {
		PrintStats(dataDir)
		return
	}

	stats := InitWordStats(
		dataDir,
		func(yield func(VocabularyPair) bool) {
			for _, lesson := range lessons {
				for _, pair := range lesson.Pairs {
					if !yield(pair) {
						return
					}
				}
			}
		},
	)

	inputReader := bufio.NewReader(os.Stdin)

	for {
		SaveStats(stats, dataDir)
		word := stats.NextWord()
		fmt.Printf("%v (%.2f): ", word.From, stats.GetStats(word).SuccessRate())
		input, _ := inputReader.ReadString('\n')
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

func PrintStats(statsDir string) {
	type KV struct {
		Key   VocabularyPair
		Value WordRecord
	}

	stats := WordStatsFromFile(statsDir).counts
	kvs := make([]KV, len(stats))
	i := 0
	for k, v := range stats {
		kvs[i] = KV{k, v}
		i++
	}
	if len(kvs) == 0 {
		return
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value.SuccessRate() > kvs[j].Value.SuccessRate()
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "From\tTo\tSuccess Rate\tTotal\n")
	for _, record := range kvs {
		fmt.Fprintf(w, "%v\t%v\t%.0f%%\t%v\n", record.Key.From, record.Key.To, record.Value.SuccessRate()*100.0, record.Value.Failure+record.Value.Success)
	}
	w.Flush()
}
