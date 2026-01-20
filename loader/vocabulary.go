package loader

import (
	"cavic19/vocabulario/core"
	"path/filepath"
	"strings"
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

	for _, row := range rows {
		if len(row) > 1 {
			newLang := strings.TrimSpace(row[0])
			// There can be multiple know words versions like new;knownA;knownB -. for each we want to create a word
			for _, knownLangRaw := range row[1:] {
				knownLang := strings.TrimSpace(knownLangRaw)
				acc[core.WordId{Word: newLang, NewLang: true}] = append(acc[core.WordId{Word: newLang, NewLang: true}], knownLang)
				acc[core.WordId{Word: knownLang, NewLang: false}] = append(acc[core.WordId{Word: knownLang, NewLang: false}], newLang)

			}

		}
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
