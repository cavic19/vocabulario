package core

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

const Rounds int = 100000

func TestVocabulary(t *testing.T) {
	testWords := [10]Word{
		testWord("A"),
		testWord("B"),
		testWord("C"),
		testWord("D"),
		testWord("E"),
		testWord("F"),
		testWord("G"),
		testWord("H"),
		testWord("I"),
		testWord("J"),
	}

	testStats := make(map[WordId]WordStat, len(testWords))

	// Initalise stats to go from the least successful to the most
	for i, word := range testWords {
		testStats[word.Id] = WordStat{
			Success: i,
			Failure: (len(testWords) - i),
		}
	}

	testVoc := NewVocabulary(testWords[:], testStats)
	freq := make(map[WordId]int, Rounds)
	// We expect the least succesfull to occur the least
	expecteSorted := make([]WordId, 0, len(testWords))
	for _, word := range testWords {
		expecteSorted = append(expecteSorted, word.Id)
	}

	for range Rounds {
		word := testVoc.NextWord()
		freq[word.Id] += 1
	}

	// Sort from most frequent to least frequent
	actualSorted := make([]WordId, 0, len(freq))
	for wordId := range freq {
		actualSorted = append(actualSorted, wordId)
	}
	sort.Slice(actualSorted, func(i, j int) bool {
		return freq[actualSorted[i]] > freq[actualSorted[j]]
	})

	assert.Equal(t, expecteSorted, actualSorted)
}

func testWord(word string) Word {
	return Word{
		Id:   WordId{Word: word, NewLang: true},
		From: word,
		To:   []string{word},
	}
}
