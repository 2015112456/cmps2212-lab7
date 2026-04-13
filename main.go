package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ---------------------------------------------------------------------------
// Domain
// ---------------------------------------------------------------------------

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ---------------------------------------------------------------------------
// Application
// ---------------------------------------------------------------------------

type application struct {
	logger *slog.Logger
	mu     sync.Mutex
	users  []User
	nextID int
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if input.Name == "" || input.Email == "" {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}

	app.mu.Lock()
	app.nextID++
	user := User{
		ID:    app.nextID,
		Name:  input.Name,
		Email: input.Email,
	}
	app.users = append(app.users, user)
	app.mu.Unlock()
	
	app.logger.Info("user created", "id", user.ID, "name", user.Name, "email", user.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (app *application) listUsers(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	users := make([]User, len(app.users))
	copy(users, app.users)
	app.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
		users:  []User{},
		nextID: 0,
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/users", app.createUser)
	mux.HandleFunc("GET /api/users", app.listUsers)

	// Static files — serves index.html, CSS, JS from ./static
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	addr := ":4000"
	logger.Info("starting server", "addr", addr)

	err := http.ListenAndServe(addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
