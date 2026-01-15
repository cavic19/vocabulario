package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

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

// Return a random word from stats based on provided statistics
// The more successful the user is with a word the less likely it is to occur
func NextWord(stats map[VocabularyPair]SuccessStats) VocabularyPair {
	N := len(stats)
	indexToPair := make([]VocabularyPair, N)
	cutOffs := make([]int, N)
	total := 0

	i := 0
	for pair, stat := range stats {
		// TODO: These should not be constants, as with more words we are getting the les likely it i
		w := int(-99*stat.SuccessRate() + 100)
		cutOffs[i] = total + w
		total += w
		indexToPair[i] = pair
		i++
	}

	r := rand.Intn(total)
	indx := rand.Intn(N)
	for i, cutOff := range cutOffs {
		if r < cutOff {
			indx = i
			break
		}
	}
	return indexToPair[indx]
}
