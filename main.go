package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"word-search-in-files/pkg/searcher"
)

func searchHandler(w http.ResponseWriter, r *http.Request, s *searcher.Searcher) {
	word := r.URL.Query().Get("word")

	files, err := s.Search(word)
	if err != nil {
		http.Error(w, "Error while searching", http.StatusInternalServerError)
		return
	}

	filesJSON, _ := json.Marshal(files)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(filesJSON)
}
func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	examplesDir := filepath.Join(currentDir, "examples")

	search := &searcher.Searcher{
		FS: os.DirFS(examplesDir),
	}

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, search)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
