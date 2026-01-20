package main

import (
	"bufio"
	"cavic19/vocabulario/core"
	"cavic19/vocabulario/loader"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
)

func main() {
	var parentDir string
	var dataDir, dataFile string
	var showStats bool

	flag.StringVar(&dataDir, "dir", "data", "Directory containing lesson CSV files")
	flag.StringVar(&dataFile, "file", "", "Specific lesson file to load (optional, loads all if not specified)")
	flag.BoolVar(&showStats, "stats", false, "Display vocabulary statistics sorted by success rate (highest to lowest)")
	flag.Parse()

	var voc *core.Vocabulary
	var err error

	if dataFile != "" {
		parentDir = filepath.Dir(dataFile)
		voc, err = loader.LoadVocabularyFromFile(dataFile)
	} else {
		parentDir = dataDir
		voc, err = loader.LoadVocabularyFromDir(dataDir)
	}

	if err != nil {
		log.Fatalf("Error loading vocabulary: %v", err)
	}

	if showStats {
		PrintStats(dataDir)
		return
	}

	// Game loop
	inputReader := bufio.NewReader(os.Stdin)
	for {
		loader.SaveStats(voc.Stats, parentDir)
		word := voc.NextWord()
		fmt.Printf("%v (%.2f): ", word.From, voc.Stats[word.Id].SuccessRate())
		input, _ := inputReader.ReadString('\n')
		if word.Test(input) {
			voc.RecordSuccess(word)
		} else {
			fmt.Printf("Wrong! It should be %v\n", word.To)
			voc.RecordFailure(word)
		}
	}
}

func PrintStats(statsDir string) {
	type KV struct {
		Key   core.WordId
		Value core.WordStat
	}

	stats := loader.LoadStats(statsDir)
	kvs := make([]KV, len(stats))
	i := 0
	for k, v := range stats {
		kvs[i] = KV{k, v}
		i++
	}
	if len(kvs) == 0 {
		return
	}

	sort.SliceStable(kvs, func(i, j int) bool {
		iRate := kvs[i].Value.SuccessRate()
		jRate := kvs[j].Value.SuccessRate()
		iTotal := kvs[i].Value.Total()
		jTotal := kvs[j].Value.Total()
		if iRate == jRate {
			if iTotal == jTotal {
				return kvs[i].Key.Word < kvs[j].Key.Word
			}
			return iTotal > jTotal
		}

		return iRate > jRate
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "From\tSuccess Rate\tTotal\n")
	for _, record := range kvs {
		fmt.Fprintf(w, "%v\t%.0f%%\t%v\n", record.Key.Word, record.Value.SuccessRate()*100.0, record.Value.Total())
	}
	w.Flush()
}
