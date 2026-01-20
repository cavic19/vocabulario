package loader

import (
	"cavic19/vocabulario/core"
	"path/filepath"
)

func LoadVocabularyFromFile(filePath string) (*core.Vocabulary, error) {
	direcotry := filepath.Dir(filePath)
	stats := LoadStats(direcotry)
	rows, err := LoadLesson(filePath)
	if err != nil {
		return nil, err
	}
	words := toWords(rows)

	return core.NewVocabulary(words, stats), nil
}

func LoadVocabularyFromDir(direcotry string) (*core.Vocabulary, error) {
	stats := LoadStats(direcotry)
	rows, err := LoadLessons(direcotry)
	if err != nil {
		return nil, err
	}
	words := toWords(rows)

	return core.NewVocabulary(words, stats), nil
}

func toWords(rows []csvRow) []core.Word {
	acc := make(map[core.WordId][]string)

	for _, pair := range rows {
		acc[core.WordId{Word: pair.First, NewLang: true}] = append(acc[core.WordId{Word: pair.First, NewLang: true}], pair.Second)
		acc[core.WordId{Word: pair.Second, NewLang: false}] = append(acc[core.WordId{Word: pair.Second, NewLang: false}], pair.First)
	}

	out := make([]core.Word, 0, len(acc))
	for k, v := range acc {
		out = append(out, core.Word{
			Id:   k,
			From: k.Word,
			To:   v,
		})
	}
	return out
}
