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

outer:
	for {
		word := words[rand.Intn(len(words))]
		for {
			fmt.Printf("%v: ", word.From)
			input, _ := reader.ReadString('\n')
			if compare(word.To, input, ComparisonOptions{}) {
				continue outer
			} else {
				fmt.Printf("Wrong! It should be %v\n", word.To)
			}
		}
	}
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
