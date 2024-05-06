package searcher

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"word-search-in-files/pkg/internal/dir"
)

type Searcher struct {
	FS fs.FS
}

// Search ищет слово в файлах и возвращает список имен файлов, содержащих это слово
func (s *Searcher) Search(word string) (files []string, err error) {
	word = strings.ReplaceAll(strings.ToLower(strings.TrimSpace(word)), " ", "")

	filesMap := make(map[string]bool)

	files, err = dir.FilesFS(s.FS, ".") // получаем список файлов
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка файлов: %v", err)
	}

	for _, file := range files {
		fileData, err := fs.ReadFile(s.FS, file) // читаем содержимое файла
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении файла %s: %v", file, err)
		}
		if strings.Contains(strings.ToLower(string(fileData)), word) {
			filesMap[file] = true
		}
	}

	result := make([]string, 0, len(filesMap))
	for file := range filesMap {
		result = append(result, file)
	}

	if len(result) == 0 {
		return nil, errors.New("слово не найдено")
	}

	return result, nil
}
