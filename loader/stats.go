package loader

import (
	"cavic19/vocabulario/core"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"strings"
)

const StatsFileName = "stats.json"

func LoadStats(dataDir string) map[core.WordId]core.WordStat {
	filePath := filepath.Join(dataDir, StatsFileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		// File doesn't exist or can't be read, return empty map
		return make(map[core.WordId]core.WordStat)
	}

	serialized := make(map[string]core.WordStat)
	err = json.Unmarshal(data, &serialized)
	if err != nil {
		log.Printf("Error unmarshalling stats: %v", err)
		return make(map[core.WordId]core.WordStat)
	}

	// Convert back to map
	stats := make(map[core.WordId]core.WordStat)
	for key, stat := range serialized {
		vocId, err := wordIdFromString(key)
		if err != nil {
			continue
		}
		stats[vocId] = stat
	}

	return stats
}

func SaveStats(stats map[core.WordId]core.WordStat, dataDir string) {
	filePath := filepath.Join(dataDir, StatsFileName)

	// We need to merge new stats into the saved stats, so we don't loose any
	// There are betetr solutions to that, but for now this si fine
	allStats := LoadStats(dataDir)
	maps.Copy(allStats, stats)

	serialized := make(map[string]core.WordStat)
	for wordId, stat := range allStats {
		if stat.Failure+stat.Success > 0 {
			key := wordIdToString(wordId)
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

func wordIdToString(key core.WordId) string {
	if key.NewLang {
		return key.Word + "_from"
	} else {
		return key.Word + "_to"
	}
}

func wordIdFromString(str string) (core.WordId, error) {
	var vocKey core.WordId
	if strings.HasSuffix(str, "_from") {
		vocKey.NewLang = true
		vocKey.Word = strings.TrimSuffix(str, "_from")
		return vocKey, nil
	} else if strings.HasSuffix(str, "_to") {
		vocKey.NewLang = false
		vocKey.Word = strings.TrimSuffix(str, "_to")
		return vocKey, nil
	} else {
		return vocKey, fmt.Errorf("%v is invlaid word id", str)
	}
}
