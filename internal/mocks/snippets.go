package mocks

import (
	"time"

	"github.com.scottyfionnghall.snippetbox/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "Test title",
	Content: "Test content...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

func (m *SnippetModel) Delete(id int) error {
	switch id {
	case 1:
		return nil
	default:
		return models.ErrNoRecord
	}
}
