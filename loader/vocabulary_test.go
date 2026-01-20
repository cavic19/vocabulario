package loader

import (
	"cavic19/vocabulario/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToWords(t *testing.T) {
	rows := []csvRow{
		{"A:0", "B:0"},
		{"A:0", "B:1"},
		{"A:1", "B:2"},
		{"A:2", "B:2"},
	}

	actual := toWords(rows)

	expected := []core.Word{
		{
			Id:   core.WordId{Word: "A:0", NewLang: true},
			From: "A:0",
			To:   []string{"B:0", "B:1"},
		},
		{
			Id:   core.WordId{Word: "A:1", NewLang: true},
			From: "A:1",
			To:   []string{"B:2"},
		},
		{
			Id:   core.WordId{Word: "A:2", NewLang: true},
			From: "A:2",
			To:   []string{"B:2"},
		},
		{
			Id:   core.WordId{Word: "B:0", NewLang: false},
			From: "B:0",
			To:   []string{"A:0"},
		},
		{
			Id:   core.WordId{Word: "B:1", NewLang: false},
			From: "B:1",
			To:   []string{"A:0"},
		},
		{
			Id:   core.WordId{Word: "B:2", NewLang: false},
			From: "B:2",
			To:   []string{"A:1", "A:2"},
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestToWords_MultipleEntries(t *testing.T) {
	rows := []csvRow{
		{"A", "B", "C"},
		{"D", "E", "C"},
	}

	actual := toWords(rows)

	expected := []core.Word{
		{
			Id:   core.WordId{Word: "A", NewLang: true},
			From: "A",
			To:   []string{"B", "C"},
		},
		{
			Id:   core.WordId{Word: "D", NewLang: true},
			From: "D",
			To:   []string{"E", "C"},
		},

		{
			Id:   core.WordId{Word: "B", NewLang: false},
			From: "B",
			To:   []string{"A"},
		},
		{
			Id:   core.WordId{Word: "C", NewLang: false},
			From: "C",
			To:   []string{"A", "D"},
		},
		{
			Id:   core.WordId{Word: "E", NewLang: false},
			From: "E",
			To:   []string{"D"},
		},
	}

	assert.ElementsMatch(t, expected, actual)

}
