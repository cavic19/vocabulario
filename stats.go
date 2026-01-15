package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

type WordRecord struct {
	Success int
	Failure int
}

func (ss WordRecord) IncrSuccess() WordRecord {
	return WordRecord{
		ss.Success + 1,
		ss.Failure,
	}
}

func (ss WordRecord) IncrFailure() WordRecord {
	return WordRecord{
		ss.Success,
		ss.Failure + 1,
	}
}

// Returns a number between 0 and 1
func (s WordRecord) SuccessRate() float32 {
	total := s.Success + s.Failure
	if total == 0 {
		return 0
	} else {
		return float32(s.Success) / float32(total)
	}
}

type WordStats struct {
	counts map[VocabularyPair]WordRecord
	memory *Memory[VocabularyPair]
}

func EmptyWordStats() *WordStats {
	return &WordStats{
		make(map[VocabularyPair]WordRecord),
		NewMemory[VocabularyPair](5),
	}
}

func (ws *WordStats) GetStats(word VocabularyPair) WordRecord {
	return ws.counts[word]
}

func (ws *WordStats) RecordFailure(word VocabularyPair) {
	old := ws.counts[word]
	ws.counts[word] = old.IncrFailure()
}

func (ws *WordStats) RecordSuccess(word VocabularyPair) {
	old := ws.counts[word]
	ws.counts[word] = old.IncrSuccess()
}

// Return a random word from stats based on provided statistics
// The more successful the user is with a word the less likely it is to occur
func (stats WordStats) NextWord() VocabularyPair {
	N := len(stats.counts)
	indexToPair := make([]VocabularyPair, N)
	cutOffs := make([]int, N)
	total := 0

	i := 0
	for pair, stat := range stats.counts {
		if stats.memory.Has(pair) {
			continue
		}
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

	nextWord := indexToPair[indx]
	stats.memory.Push(nextWord)
	return nextWord
}

func SaveStats(stats *WordStats, dataDir string) {
	filePath := filepath.Join(dataDir, "stats.json")

	serialized := make(map[string]WordRecord)
	for pair, stat := range stats.counts {
		if stat.Failure+stat.Success > 0 {
			key := pair.From + "->" + pair.To
			serialized[key] = stat
		}
	}

	if len(serialized) == 0 {
		return
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

func LoadStats(dataDir string) map[VocabularyPair]WordRecord {
	filePath := filepath.Join(dataDir, "stats.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		// File doesn't exist or can't be read, return empty map
		return make(map[VocabularyPair]WordRecord)
	}

	serialized := make(map[string]WordRecord)
	err = json.Unmarshal(data, &serialized)
	if err != nil {
		log.Printf("Error unmarshalling stats: %v", err)
		return make(map[VocabularyPair]WordRecord)
	}

	// Convert back to map
	stats := make(map[VocabularyPair]WordRecord)
	for key, stat := range serialized {
		parts := strings.Split(key, "->")
		if len(parts) == 2 {
			pair := VocabularyPair{From: parts[0], To: parts[1]}
			stats[pair] = stat
		}
	}

	return stats
}
