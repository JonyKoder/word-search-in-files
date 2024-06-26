package searcher

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestSearcher_Search_NotFound(t *testing.T) {
	s := &Searcher{
		FS: fstest.MapFS{
			"file1.txt": {Data: []byte("Hello")},
			"file2.txt": {Data: []byte("World")},
		},
	}
	word := "NotInFiles"

	gotFiles, err := s.Search(word)

	if err != nil {
		t.Errorf("Search(%s) вернул ошибку: %v", word, err)
	}

	if gotFiles != nil {
		t.Errorf("Search(%s) вернул файлы, когда их не ожидалось: %v", word, gotFiles)
	}
}

func TestSearcher_Search_ReadFileError(t *testing.T) {
	s := &Searcher{
		FS: fstest.MapFS{
			"file1.txt": {Data: []byte("Hello")},
			"file2.txt": {Data: []byte("World")},
		},
	}

	// Создаем временный каталог для тестовых файлов
	tempDir := t.TempDir()
	word := "дети"

	// Путь к файлу
	filePath := filepath.Join(tempDir, "file1.txt")

	// Записываем данные в файл
	if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
		t.Fatalf("ошибка при создании файла: %v", err)
	}

	// Вызываем функцию поиска
	_, err := s.Search(word)

	// Проверяем, что вернулась ошибка
	if err == nil {
		t.Errorf("ожидалась ошибка при чтении файлов, но вернулось nil")
	}
}

func TestSearcher_Search(t *testing.T) {
	type fields struct {
		FS fs.FS
	}
	type args struct {
		word string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantFiles []string
		wantErr   bool
	}{
		{
			name: "Ok",
			fields: fields{
				FS: fstest.MapFS{
					"file1.txt": {Data: []byte("World")},
					"file2.txt": {Data: []byte("World1")},
					"file3.txt": {Data: []byte("Hello World")},
				},
			},
			args:      args{word: "World"},
			wantFiles: []string{"file1", "file3"},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Searcher{
				FS: tt.fields.FS,
			}
			gotFiles, err := s.Search(tt.args.word)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("Search() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}
