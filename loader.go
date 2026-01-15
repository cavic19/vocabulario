package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func loadLessons(dataDir string, lessonName string) ([]Lesson, error) {
	var lessons []Lesson

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".csv") {
			name := strings.TrimSuffix(entry.Name(), ".csv")

			// Filter by lesson name if specified
			if lessonName != "" && name != lessonName {
				continue
			}

			lesson, err := loadLesson(filepath.Join(dataDir, entry.Name()), entry.Name())
			if err != nil {
				return nil, err
			}
			lessons = append(lessons, lesson)
		}
	}

	return lessons, nil
}

func loadLesson(filePath, fileName string) (Lesson, error) {
	lesson := Lesson{
		Name: strings.TrimSuffix(fileName, ".csv"),
	}

	file, err := os.Open(filePath)
	if err != nil {
		return lesson, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ";")
		if len(parts) == 2 {
			pair := VocabularyPair{
				From: strings.TrimSpace(parts[0]),
				To:   strings.TrimSpace(parts[1]),
			}
			lesson.Pairs = append(lesson.Pairs, pair)
		}
	}

	if err := scanner.Err(); err != nil {
		return lesson, err
	}

	return lesson, nil
}
