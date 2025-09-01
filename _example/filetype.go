package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func DetectContentType(data []byte) string {
	return http.DetectContentType(data)
}

func main() {
	files_name := []string{"./_example/sample1.bmp", "./_example/sample1.png", "./_example/fake_sample2.jpg", "./_example/avatar.jpg", "./_example/github-card.png"}
	var files_path []string
	for _, f := range files_name {
		file_path, err := filepath.Abs(f)
		if err != nil {
			continue
		}
		files_path = append(files_path, file_path)
	}

	for _, f := range files_path {
		file, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		mime_type := DetectContentType(file)
		fmt.Printf("the file %v mimetype is %v\n", f, mime_type)
	}
}
