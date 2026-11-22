package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	s, err := openFileStore(filepath.Join(t.TempDir(), "tasks.json"))
	if err != nil {
		t.Fatal(err)
	}
	return newHandler(s)
}
func request(t *testing.T, h http.Handler, method, p string, body string) (int, []byte) {
	t.Helper()
	r := httptest.NewRequest(method, p, strings.NewReader(body))
	if body != "" {
		r.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w.Code, w.Body.Bytes()
}
func TestCRUD(t *testing.T) {
	h := newTestServer(t)
	code, b := request(t, h, "GET", "/api/tasks", "")
	if code != 200 || string(bytes.TrimSpace(b)) != "[]" {
		t.Fatalf("empty list: %d %s", code, b)
	}
	code, b = request(t, h, "POST", "/api/tasks", `{"title":"Test task","description":"hello","dueDate":"2026-11-30","priority":"high"}`)
	if code != 201 {
		t.Fatalf("create %d %s", code, b)
	}
	var task Task
	if err := json.Unmarshal(b, &task); err != nil {
		t.Fatal(err)
	}
	if task.ID != 1 || task.DueDate != "2026-11-30" {
		t.Fatal(task)
	}
	code, b = request(t, h, "GET", "/api/tasks/1", "")
	if code != 200 {
		t.Fatalf("get: %d %s", code, b)
	}
	code, b = request(t, h, "PUT", "/api/tasks/1", `{"title":"Updated","status":"done","priority":"low"}`)
	if code != 200 {
		t.Fatalf("put: %d %s", code, b)
	}
	json.Unmarshal(b, &task)
	if !task.Done || task.Status != "done" {
		t.Fatal(task)
	}
	code, b = request(t, h, "GET", "/api/stats", "")
	if code != 200 || !bytes.Contains(b, []byte(`"done":1`)) {
		t.Fatalf("stats %d %s", code, b)
	}
	code, b = request(t, h, "GET", "/api/tasks?q=updated", "")
	if code != 200 || !bytes.Contains(b, []byte("Updated")) {
		t.Fatalf("search: %d %s", code, b)
	}
	code, b = request(t, h, "GET", "/api/export.csv", "")
	if code != 200 || !bytes.Contains(b, []byte("Updated")) {
		t.Fatalf("csv %d %s", code, b)
	}
	code, b = request(t, h, "DELETE", "/api/tasks/1", "")
	if code != 204 {
		t.Fatalf("delete %d %s", code, b)
	}
	code, b = request(t, h, "GET", "/api/tasks/1", "")
	if code != 404 {
		t.Fatalf("missing %d %s", code, b)
	}
}
func TestValidationAndStatic(t *testing.T) {
	h := newTestServer(t)
	for _, v := range []string{`{"title":""}`, `{"title":"ok","status":"wrong"}`, `{"title":"ok","due_date":"next week"}`, `{"title":"ok","priority":"urgent"}`, `[]`} {
		code, b := request(t, h, "POST", "/api/tasks", v)
		if code != 400 {
			t.Fatalf("accepted invalid %s: %d %s", v, code, b)
		}
	}
	for _, p := range []string{"/", "/board", "/vendor/react.production.min.js", "/app.js", "/style.css"} {
		code, b := request(t, h, "GET", p, "")
		if code != 200 || len(b) < 10 {
			t.Fatalf("static %s %d %q", p, code, b)
		}
	}
	code, _ := request(t, h, "GET", "/missing.js", "")
	if code != 404 {
		t.Fatalf("missing.js: %d", code)
	}
}
func TestPersistence(t *testing.T) {
	p := filepath.Join(t.TempDir(), "tasks.json")
	s, _ := openFileStore(p)
	task, _ := normalizeInput(TaskInput{Title: "Stored"}, nil)
	if _, err := s.Create(task); err != nil {
		t.Fatal(err)
	}
	newS, err := openFileStore(p)
	if err != nil {
		t.Fatal(err)
	}
	items, err := newS.List()
	if err != nil || len(items) != 1 || items[0].Title != "Stored" {
		t.Fatalf("persistence: %v %v", items, err)
	}
}
func TestCors(t *testing.T) {
	h := newTestServer(t)
	r := httptest.NewRequest("OPTIONS", "/api/tasks", nil)
	r.Header.Set("Origin", "https://evil.example")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal(w.Code)
	}
	r = httptest.NewRequest("OPTIONS", "/api/tasks", nil)
	r.Header.Set("Origin", "http://localhost:5173")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 204 {
		t.Fatal(w.Code)
	}
}
