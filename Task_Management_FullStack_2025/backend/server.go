package main

import (
	"embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//go:embed web
var webFiles embed.FS

func main() {
	s, err := newStore()
	if err != nil {
		log.Fatal("Storage startup: ", err)
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "127.0.0.1:8082"
	}
	log.Printf("Task Manager Full-Stack started | http://localhost:8082 | storage initialized (%s)", addr)
	log.Fatal(http.ListenAndServe(addr, newHandler(s)))
}
func newHandler(s Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			notAllowed(w)
			return
		}
		writeJSON(w, 200, map[string]interface{}{"ok": true, "app": "Task Management Full-Stack", "timestamp": time.Now().UTC().Format(time.RFC3339)})
	})
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tasks, err := s.List()
			if err != nil {
				serverError(w, err)
				return
			}
			q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
			status := r.URL.Query().Get("status")
			priority := r.URL.Query().Get("priority")
			out := make([]Task, 0, len(tasks))
			for _, t := range tasks {
				if q != "" && !strings.Contains(strings.ToLower(t.Title+" "+t.Description), q) {
					continue
				}
				if status != "" && t.Status != status {
					continue
				}
				if priority != "" && t.Priority != priority {
					continue
				}
				out = append(out, t)
			}
			writeJSON(w, 200, out)
		case "POST":
			input, err := readInput(r)
			if err != nil {
				badRequest(w, err.Error())
				return
			}
			task, err := normalizeInput(input, nil)
			if err != nil {
				badRequest(w, err.Error())
				return
			}
			created, err := s.Create(task)
			if err != nil {
				serverError(w, err)
				return
			}
			w.Header().Set("Location", fmt.Sprintf("/api/tasks/%d", created.ID))
			writeJSON(w, 201, created)
		default:
			notAllowed(w)
		}
	})
	mux.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		idstr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
		if strings.Contains(idstr, "/") || idstr == "" {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil || id < 1 {
			badRequest(w, "Неверный ID задачи")
			return
		}
		switch r.Method {
		case "GET":
			t, ok, err := s.Get(id)
			if err != nil {
				serverError(w, err)
				return
			}
			if !ok {
				notFound(w)
				return
			}
			writeJSON(w, 200, t)
		case "PUT", "PATCH":
			old, ok, err := s.Get(id)
			if err != nil {
				serverError(w, err)
				return
			}
			if !ok {
				notFound(w)
				return
			}
			input, err := readInput(r)
			if err != nil {
				badRequest(w, err.Error())
				return
			}
			t, err := normalizeInput(input, &old)
			if err != nil {
				badRequest(w, err.Error())
				return
			}
			saved, ok, err := s.Update(id, t)
			if err != nil {
				serverError(w, err)
				return
			}
			if !ok {
				notFound(w)
				return
			}
			writeJSON(w, 200, saved)
		case "DELETE":
			ok, err := s.Delete(id)
			if err != nil {
				serverError(w, err)
				return
			}
			if !ok {
				notFound(w)
				return
			}
			w.WriteHeader(204)
		default:
			notAllowed(w)
		}
	})
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			notAllowed(w)
			return
		}
		tasks, err := s.List()
		if err != nil {
			serverError(w, err)
			return
		}
		counts := map[string]int{"total": len(tasks), "todo": 0, "in_progress": 0, "done": 0, "overdue": 0}
		today := time.Now().Format("2006-01-02")
		for _, t := range tasks {
			counts[t.Status]++
			if t.Status != "done" && t.DueDate != "" && t.DueDate < today {
				counts["overdue"]++
			}
		}
		writeJSON(w, 200, counts)
	})
	mux.HandleFunc("/api/export.csv", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			notAllowed(w)
			return
		}
		tasks, err := s.List()
		if err != nil {
			serverError(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=tasks.csv")
		w.Write([]byte{0xEF, 0xBB, 0xBF})
		c := csv.NewWriter(w)
		_ = c.Write([]string{"ID", "Title", "Description", "Status", "Priority", "Due date", "Created"})
		for _, t := range tasks {
			_ = c.Write([]string{strconv.FormatInt(t.ID, 10), t.Title, t.Description, t.Status, t.Priority, t.DueDate, t.CreatedAt.Format(time.RFC3339)})
		}
		c.Flush()
	})
	static, err := fs.Sub(webFiles, "web")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(static))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			notAllowed(w)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			notFound(w)
			return
		}
		if r.URL.Path != "/" {
			p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
			if file, err := fs.Stat(static, p); err == nil && !file.IsDir() {
				fileServer.ServeHTTP(w, r)
				return
			}
			if strings.Contains(path.Base(p), ".") {
				notFound(w)
				return
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		b, err := fs.ReadFile(static, "index.html")
		if err != nil {
			serverError(w, err)
			return
		}
		w.Write(b)
	})
	return corsSecurity(mux)
}
func corsSecurity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		origin := r.Header.Get("Origin")
		if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" && origin != "" {
			sameOrigin := origin == "http://"+r.Host || origin == "https://"+r.Host
			devOrigin := origin == "http://localhost:3000" || origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173"
			if !sameOrigin && !devOrigin {
				writeJSON(w, 403, map[string]string{"error": "Запрещённый Origin"})
				return
			}
		}
		if origin == "http://localhost:3000" || origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		}
		if r.Method == "OPTIONS" {
			if origin != "" && w.Header().Get("Access-Control-Allow-Origin") == "" {
				w.WriteHeader(403)
			} else {
				w.WriteHeader(204)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
func readInput(r *http.Request) (TaskInput, error) {
	defer r.Body.Close()
	var raw map[string]json.RawMessage
	dec := json.NewDecoder(io.LimitReader(r.Body, 1024*1024))
	if err := dec.Decode(&raw); err != nil {
		return TaskInput{}, fmt.Errorf("Требуется JSON-объект: %w", err)
	}
	if raw == nil {
		return TaskInput{}, errors.New("Требуется JSON-объект")
	}
	var extra interface{}
	if err := dec.Decode(&extra); err != io.EOF {
		return TaskInput{}, errors.New("Допускается только один JSON-объект")
	}
	if _, ok := raw["due_date"]; !ok {
		if due, ok := raw["dueDate"]; ok {
			raw["due_date"] = due
		}
	}
	b, _ := json.Marshal(raw)
	var input TaskInput
	if err := json.Unmarshal(b, &input); err != nil {
		return TaskInput{}, err
	}
	return input, nil
}
func writeJSON(w http.ResponseWriter, code int, val interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(val)
}
func badRequest(w http.ResponseWriter, s string) { writeJSON(w, 400, map[string]string{"error": s}) }
func notFound(w http.ResponseWriter) {
	writeJSON(w, 404, map[string]string{"error": "Задача не найдена"})
}
func notAllowed(w http.ResponseWriter) {
	writeJSON(w, 405, map[string]string{"error": "Метод не поддерживается"})
}
func serverError(w http.ResponseWriter, err error) {
	log.Print("Server error: ", err)
	writeJSON(w, 500, map[string]string{"error": "Ошибка сохранения данных"})
}
