package docs

import (
	"os"
	"path/filepath"
	"strings"
)

type Document struct {
	Name    string
	Path    string
	Content string
}

type DocStore struct {
	docs []Document
}

func NewDocStore() *DocStore {
	return &DocStore{
		docs: make([]Document, 0),
	}
}

func (s *DocStore) LoadFromDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(filepath.Base(file), ".md")
		s.docs = append(s.docs, Document{
			Name:    name,
			Path:    file,
			Content: string(content),
		})
	}

	return nil
}

func (s *DocStore) GetDocs() []Document {
	return s.docs
}

func (s *DocStore) GetDoc(name string) *Document {
	for i := range s.docs {
		if s.docs[i].Name == name {
			return &s.docs[i]
		}
	}
	return nil
}
