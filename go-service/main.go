package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Item struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Store struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]Item
}

func NewStore() *Store {
	return &Store{items: make(map[int64]Item)}
}

func (s *Store) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Item, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}
	return result
}

func (s *Store) Create(req ItemRequest) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	item := Item{ID: s.nextID, Name: req.Name, Description: req.Description}
	s.items[item.ID] = item
	return item
}

func (s *Store) Get(id int64) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	return item, ok
}

func (s *Store) Update(id int64, req ItemRequest) (Item, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[id]
	if !ok {
		return Item{}, false
	}
	item := Item{ID: id, Name: req.Name, Description: req.Description}
	s.items[id] = item
	return item, true
}

func (s *Store) Delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}

func main() {
	store := NewStore()

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, store.List())
		case http.MethodPost:
			var req ItemRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json body", http.StatusBadRequest)
				return
			}
			item := store.Create(req)
			writeJSON(w, http.StatusCreated, item)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		idPart := strings.TrimPrefix(r.URL.Path, "/items/")
		id, err := strconv.ParseInt(idPart, 10, 64)
		if err != nil {
			http.Error(w, "invalid item id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			item, ok := store.Get(id)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, item)
		case http.MethodPut:
			var req ItemRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json body", http.StatusBadRequest)
				return
			}
			item, ok := store.Update(id, req)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, item)
		case http.MethodDelete:
			if !store.Delete(id) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Go service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
