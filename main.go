package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"strconv"
)

// ---------------------------------------------------------------------------
// Domain
// ---------------------------------------------------------------------------

type Event struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Tickets string `json:"tickets"`
	Terms bool   `json:"terms"`
}

// ---------------------------------------------------------------------------
// Application
// ---------------------------------------------------------------------------

type application struct {
	logger *slog.Logger
	mu     sync.Mutex
	events  []Event
	nextID int
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func (app *application) registerEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Date  string `json:"date"`
		Tickets string `json:"tickets"`
		Terms bool   `json:"terms"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	input.Date = strings.TrimSpace(input.Date)
	input.Tickets = strings.TrimSpace(input.Tickets)

	if input.Date == "" || input.Tickets == "" {
		http.Error(w, "Date and Number of Tickets are required", http.StatusBadRequest)
		return
	}
	if !input.Terms {
		http.Error(w, "You must agree to the Terms and Conditions", http.StatusBadRequest)
		return
	}

	
	// Validate date format and ensure it's in the future
	today := time.Now().Truncate(24 * time.Hour)
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		http.Error(w, "Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}
	if date.Before(today) {
		http.Error(w, "Date must be today or in the future", http.StatusBadRequest)
		return
	}

	// Validate ticket range
	tickets, err := strconv.Atoi(input.Tickets)
	if err != nil {
		http.Error(w, "Please enter a valid integer.", http.StatusBadRequest)
		return
	}
	if tickets < 1 || tickets > 5 {
		http.Error(w, "Invalid number of tickets. Tickets must be between 1 and 5", http.StatusBadRequest)
		return
	}

	app.mu.Lock()
	app.nextID++
	event := Event{
		ID:    app.nextID,
		Date:  input.Date,
		Tickets: input.Tickets,
		Terms: input.Terms,
	}
	app.events = append(app.events, event)
	app.mu.Unlock()
	
	app.logger.Info("Event registered", "ID", event.ID, "Date", event.Date, "Tickets", event.Tickets)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

func (app *application) listEventRegistrations(w http.ResponseWriter, r *http.Request) {
	app.mu.Lock()
	events := make([]Event, len(app.events))
	copy(events, app.events)
	app.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
		events:  []Event{},
		nextID: 0,
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			app.registerEvent(w, r)
		case http.MethodGet:
			app.listEventRegistrations(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Static files — serves index.html, CSS, JS from ./static
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	addr := ":4000"
	logger.Info("starting server", "addr", addr)

	err := http.ListenAndServe(addr, mux)
	logger.Error(err.Error())
}
