package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

	lessons, err := loadLessons(dataDir, lessonName)
	if err != nil {
		log.Fatalf("Error loading lessons: %v", err)
	}

	if len(lessons) == 0 {
		fmt.Println("No lessons found.")
		return
	}

	var words []VocabularyPair

	for _, lesson := range lessons {
		for _, pair := range lesson.Pairs {
			words = append(words, VocabularyPair{pair.From, pair.To}, VocabularyPair{pair.To, pair.From})
		}
	}
	reader := bufio.NewReader(os.Stdin)
	var stats []SuccessStats = make([]SuccessStats, len(words))

	for {
		indx := NextIdx(stats)
		word := words[indx]

		fmt.Printf("%v: ", word.From)
		input, _ := reader.ReadString('\n')
		oldStats := stats[indx]
		if compare(word.To, input, ComparisonOptions{}) {
			stats[indx] = SuccessStats{
				oldStats.Success + 1,
				oldStats.Failure,
			}
		} else {
			fmt.Printf("Wrong! It should be %v\n", word.To)
			stats[indx] = SuccessStats{
				oldStats.Success,
				oldStats.Failure + 1,
			}
		}
	}
}

// index -> to stat
func NextIdx(statsByIdx []SuccessStats) int {
	cutOffs := make([]int, len(statsByIdx))
	total := 0
	for i, stat := range statsByIdx {
		w := int(-9*stat.SuccessRate() + 10)
		cutOffs[i] = total + w
		total += w
	}

	r := rand.Intn(total)
	indx := rand.Intn(len(statsByIdx))
	for i, cutOff := range cutOffs {
		if r < cutOff {
			indx = i
			break
		}
	}
	return indx
}

type ComparisonResult struct {
}

type ComparisonOptions struct {
	ignoreAccents bool
	ignoreCases   bool
}

func compare(expected, actual string, opts ComparisonOptions) bool {
	e, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(expected)))
	a, _, _ := transform.String(transformer, strings.ToLower(strings.TrimSpace(actual)))
	return e == a
}

func loadLessons(dataDir string, lessonName string) ([]Lesson, error) {
	var lessons []Lesson

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".csv") {
			name := strings.TrimSuffix(entry.Name(), ".csv")

			// Filter by lesson name if specified
			if lessonName != "" && name != lessonName {
				continue
			}

			lesson, err := loadLesson(filepath.Join(dataDir, entry.Name()), entry.Name())
			if err != nil {
				return nil, err
			}
			lessons = append(lessons, lesson)
		}
	}

	return lessons, nil
}

func loadLesson(filePath, fileName string) (Lesson, error) {
	lesson := Lesson{
		Name: strings.TrimSuffix(fileName, ".csv"),
	}

	file, err := os.Open(filePath)
	if err != nil {
		return lesson, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ";")
		if len(parts) == 2 {
			pair := VocabularyPair{
				From: strings.TrimSpace(parts[0]),
				To:   strings.TrimSpace(parts[1]),
			}
			lesson.Pairs = append(lesson.Pairs, pair)
		}
	}

	if err := scanner.Err(); err != nil {
		return lesson, err
	}

	return lesson, nil
}

type SuccessStats struct {
	Success int
	Failure int
}

// Return number between 0 and 1
func (s SuccessStats) SuccessRate() float32 {
	total := s.Success + s.Failure
	if total == 0 {
		return 0
	} else {
		return float32(s.Success) / float32(total)
	}
}

// // I don't want to repeat a word for say 3 steps?
// // I want to success rate to dictate probability
// // I want success rate to update as I keep guessing and occasionally I want to flush it to a file, flushing happens also on keybaord interrupt and such
// func NextWord(words []VocabularyPair, rates map[string]SuccessRate) VocabularyPair {

// }
