package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
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
		if compare(word.To, input, ComparisonOptions{}) {
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

func SaveStats(stats map[VocabularyPair]SuccessStats, dataDir string) {
	filePath := filepath.Join(dataDir, "stats.json")

	serialized := make(map[string]SuccessStats)
	for pair, stat := range stats {
		if stat.Failure+stat.Success > 0 {
			key := pair.From + "->" + pair.To
			serialized[key] = stat
		}
	}

	data, err := json.MarshalIndent(serialized, "", "  ")
	if err != nil {
		log.Printf("Error marshalling stats: %v", err)
		return
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		log.Printf("Error writing stats file: %v", err)
		return
	}
}

func LoadStats(dataDir string) map[VocabularyPair]SuccessStats {
	filePath := filepath.Join(dataDir, "stats.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		// File doesn't exist or can't be read, return empty map
		return make(map[VocabularyPair]SuccessStats)
	}

	serialized := make(map[string]SuccessStats)
	err = json.Unmarshal(data, &serialized)
	if err != nil {
		log.Printf("Error unmarshalling stats: %v", err)
		return make(map[VocabularyPair]SuccessStats)
	}

	// Convert back to map
	stats := make(map[VocabularyPair]SuccessStats)
	for key, stat := range serialized {
		parts := strings.Split(key, "->")
		if len(parts) == 2 {
			pair := VocabularyPair{From: parts[0], To: parts[1]}
			stats[pair] = stat
		}
	}

	return stats
}

func NextWord(stats map[VocabularyPair]SuccessStats) VocabularyPair {
	indexToPair := make([]VocabularyPair, len(stats))
	cutOffs := make([]int, len(stats))
	total := 0

	i := 0
	for pair, stat := range stats {
		w := int(-9*stat.SuccessRate() + 10)
		cutOffs[i] = total + w
		total += w
		indexToPair[i] = pair
		i++
	}

	r := rand.Intn(total)
	indx := rand.Intn(len(stats))
	for i, cutOff := range cutOffs {
		if r < cutOff {
			indx = i
			break
		}
	}
	return indexToPair[indx]
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
