package app

import (
	"encoding/json"
	"fmt"
	"gredis/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct {
	Articles map[string]db.Article
}

func (m *mockDB) GetArticleById(id string) db.Article {
	return m.Articles[id]
}

func TestGetArticlesById(t *testing.T) {
	m := &mockDB{
		Articles: map[string]db.Article{
			"67": {
				ArticleID: 67,
				Title:     "cf104098f7e25daf58b87a5f2a4fd7b8",
				Text:      "026ae4494731b72ea06795a3e3f489fb",
			},
		},
	}
	// Create a new HTTP request with a placeholder for the ID value
	r := httptest.NewRequest("GET", "/articles/{id}", nil)

	// Set the value of the ID path parameter
	fmt.Println("pathvalue", r.PathValue("id"))
	w := httptest.NewRecorder()
	getArticlesById(m)(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	b, err := json.Marshal(m.Articles["67"])
	if err != nil {
		t.Error(err)
	}
	if string(b) != string(w.Body.Bytes()) {
		t.Errorf("expected body %s, got %s", string(b), string(w.Body.Bytes()))
	}
}

func TestWelcome(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	welcome(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
