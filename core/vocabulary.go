package core

import (
	"math/rand"
)

type Vocabulary struct {
	Stats  map[WordId]WordStat
	Words  map[WordId]Word
	memory *Memory[WordId]
}

const WordMemorySize int = 5

func NewVocabulary(words []Word, stats map[WordId]WordStat) *Vocabulary {
	actualStats := make(map[WordId]WordStat, len(words))
	wordLookup := make(map[WordId]Word, len(words))
	for _, word := range words {
		wordLookup[word.Id] = word
		actualStats[word.Id] = stats[word.Id]
	}
	return &Vocabulary{
		Stats:  actualStats,
		Words:  wordLookup,
		memory: NewMemory[WordId](WordMemorySize),
	}
}

func EmptyVocabulary() *Vocabulary {
	return NewVocabulary(
		[]Word{},
		make(map[WordId]WordStat),
	)
}

func (ws *Vocabulary) RecordFailure(word Word) {
	old := ws.Stats[word.Id]
	ws.Stats[word.Id] = old.IncrFailure()
}

func (ws *Vocabulary) RecordSuccess(word Word) {
	old := ws.Stats[word.Id]
	ws.Stats[word.Id] = old.IncrSuccess()
}

// Return a random word from stats based on provided statistics
// The more successful the user is with a word the less likely it is to occur
func (stats Vocabulary) NextWord() Word {
	N := len(stats.Stats)
	indexToPair := make([]Word, N)
	// Helper
	cutOffs := make([]int, N)
	total := 0

	i := 0
	for pairKey, stat := range stats.Stats {
		if stats.memory.Has(pairKey) {
			continue
		}
		// 100% success rate -> 1
		// 0% success rate -> N
		w := int(float32(1-N)*stat.SuccessRate() + float32(N))
		cutOffs[i] = total + w
		total += w
		indexToPair[i] = stats.Words[pairKey]
		i++
	}

	r := rand.Intn(total)
	indx := rand.Intn(len(indexToPair))
	for i, cutOff := range cutOffs {
		if r < cutOff {
			indx = i
			break
		}
	}

	nextWord := indexToPair[indx]
	stats.memory.Push(nextWord.Id)
	return nextWord
}
