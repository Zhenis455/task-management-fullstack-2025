//go:build postgres

package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

// Production / classroom PostgreSQL adapter using GORM.
// Build with: go run -tags postgres .
type postgresStore struct{ db *gorm.DB }

func newStore() (Store, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required for PostgreSQL mode")
	}
	var db *gorm.DB
	var err error
	for attempt := 0; attempt < 25; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := db.AutoMigrate(&Task{}); err != nil {
		return nil, err
	}
	return &postgresStore{db: db}, nil
}
func (s *postgresStore) List() ([]Task, error) {
	var tasks []Task
	err := s.db.Order("id DESC").Find(&tasks).Error
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, err
}
func (s *postgresStore) Get(id int64) (Task, bool, error) {
	var t Task
	err := s.db.First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Task{}, false, nil
	}
	return t, err == nil, err
}
func (s *postgresStore) Create(t Task) (Task, error) {
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = t.CreatedAt
	err := s.db.Create(&t).Error
	return t, err
}
func (s *postgresStore) Update(id int64, t Task) (Task, bool, error) {
	old, ok, err := s.Get(id)
	if err != nil || !ok {
		return Task{}, ok, err
	}
	t.ID = id
	t.CreatedAt = old.CreatedAt
	t.UpdatedAt = time.Now().UTC()
	err = s.db.Model(&Task{}).Where("id = ?", id).Select("Title", "Description", "Done", "DueDate", "Priority", "Status", "UpdatedAt").Updates(&t).Error
	return t, err == nil, err
}
func (s *postgresStore) Delete(id int64) (bool, error) {
	res := s.db.Delete(&Task{}, id)
	return res.RowsAffected > 0, res.Error
}
