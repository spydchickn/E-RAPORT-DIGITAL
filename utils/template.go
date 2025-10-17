package utils

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Templates *template.Template

func loadTemplates() *template.Template {
	t := template.New("")
	err := filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			relPath, err := filepath.Rel("views", path)
			if err != nil {
				return err
			}
			relPath = strings.Replace(relPath, "\\", "/", -1) // Normalize for Windows
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			t = t.New(relPath)
			_, err = t.Parse(string(content))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error loading templates:", err)
	}
	return t
}

func init() {
	Templates = loadTemplates()
}

// RenderPartial executes a named template to a string for injection into base layout
func RenderPartial(name string, data interface{}) string {
	var b strings.Builder
	err := Templates.ExecuteTemplate(&b, name, data)
	if err != nil {
		log.Printf("Error rendering partial %s: %v", name, err)
		return "" // or return error, but for simplicity, empty string
	}
	return b.String()
}
