package loader

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type csvRow []string

func LoadLessons(dataDir string) ([]csvRow, error) {
	var rows []csvRow

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".csv") {
			lessonRows, err := LoadLesson(filepath.Join(dataDir, entry.Name()))
			if err != nil {
				return nil, err
			}
			rows = append(rows, lessonRows...)
		}
	}

	return rows, nil
}

func LoadLesson(filePath string) ([]csvRow, error) {
	rows := make([]csvRow, 0)
	file, err := os.Open(filePath)
	if err != nil {
		return rows, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ";")
		rows = append(rows, parts)
	}

	if err := scanner.Err(); err != nil {
		return rows, err
	}

	return rows, nil
}
