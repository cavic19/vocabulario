package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	assert.True(t, compare("España", "espana"))
}

func TestWordTest(t *testing.T) {
	word := Word{
		Id:   WordId{Word: "pan", NewLang: true},
		From: "pan",
		To:   []string{"bread", "loaf", "pan"},
	}

	assert.True(t, word.Test("BREAD"))
	assert.True(t, word.Test("Loaf"))
	assert.True(t, word.Test("pan"))
}
