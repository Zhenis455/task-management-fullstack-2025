//go:build !postgres

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type fileDB struct {
	NextID int64  `json:"next_id"`
	Tasks  []Task `json:"tasks"`
}
type fileStore struct {
	mu   sync.Mutex
	path string
	db   fileDB
}

func newStore() (Store, error) {
	path := os.Getenv("TASKS_FILE")
	if path == "" {
		path = filepath.Join("data", "tasks.json")
	}
	return openFileStore(path)
}
func openFileStore(path string) (Store, error) {
	s := &fileStore{path: path, db: fileDB{NextID: 1, Tasks: []Task{}}}
	b, err := os.ReadFile(path)
	if err == nil {
		if err = json.Unmarshal(b, &s.db); err != nil {
			return nil, fmt.Errorf("повреждён JSON-файл %s: %w", path, err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if s.db.NextID < 1 {
		s.db.NextID = 1
	}
	for _, t := range s.db.Tasks {
		if t.ID >= s.db.NextID {
			s.db.NextID = t.ID + 1
		}
	}
	return s, nil
}
func (s *fileStore) save(next fileDB) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err = os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	if err = os.Rename(tmp, s.path); err != nil {
		os.Remove(tmp)
		return err
	}
	s.db = next
	return nil
}
func (s *fileStore) List() ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := append([]Task{}, s.db.Tasks...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	return out, nil
}
func (s *fileStore) Get(id int64) (Task, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.db.Tasks {
		if t.ID == id {
			return t, true, nil
		}
	}
	return Task{}, false, nil
}
func (s *fileStore) Create(t Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := fileDB{NextID: s.db.NextID + 1, Tasks: append([]Task{}, s.db.Tasks...)}
	t.ID = s.db.NextID
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = t.CreatedAt
	next.Tasks = append(next.Tasks, t)
	if err := s.save(next); err != nil {
		return Task{}, err
	}
	return t, nil
}
func (s *fileStore) Update(id int64, t Task) (Task, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := fileDB{NextID: s.db.NextID, Tasks: append([]Task{}, s.db.Tasks...)}
	for i, old := range next.Tasks {
		if old.ID == id {
			t.ID = id
			t.CreatedAt = old.CreatedAt
			t.UpdatedAt = time.Now().UTC()
			next.Tasks[i] = t
			if err := s.save(next); err != nil {
				return Task{}, false, err
			}
			return t, true, nil
		}
	}
	return Task{}, false, nil
}
func (s *fileStore) Delete(id int64) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	next := fileDB{NextID: s.db.NextID, Tasks: make([]Task, 0, len(s.db.Tasks))}
	found := false
	for _, t := range s.db.Tasks {
		if t.ID == id {
			found = true
		} else {
			next.Tasks = append(next.Tasks, t)
		}
	}
	if !found {
		return false, nil
	}
	if err := s.save(next); err != nil {
		return false, err
	}
	return true, nil
}
