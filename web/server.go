package web

// Web UI server for the Log Aggregator Dashboard.
// Now with bcrypt-based authentication to protect sensitive log data.
// Uses only net/http from the standard library (no external web frameworks).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"logmerge/internal/auth"
	"logmerge/internal/filter"
	"logmerge/internal/models"
	"logmerge/internal/ordering"
	"logmerge/internal/parser"
	"logmerge/internal/source"
)

// APIResponse wraps API responses with metadata.
type APIResponse struct {
	Total  int               `json:"total"`
	Events []models.LogEvent `json:"events"`
	Stats  *StatsResponse    `json:"stats,omitempty"`
}

// StatsResponse contains aggregated statistics.
type StatsResponse struct {
	TotalEvents int            `json:"total_events"`
	LevelCounts map[string]int `json:"level_counts"`
	Sources     []string       `json:"sources"`
}

// Server holds the web server state.
type Server struct {
	files    []string
	storage  *models.LogStorage
	events   []models.LogEvent
	users    *auth.UserStore      // BCRYPT: user credential store
	sessions *auth.SessionManager // Session token management
	password string               // dashboard password (empty = no auth)
}

// NewServer creates a new web server for the given log files.
// If password is non-empty, the dashboard requires authentication.
func NewServer(files []string, password string) *Server {
	s := &Server{
		files:    files,
		storage:  models.NewLogStorage(),
		events:   make([]models.LogEvent, 0),
		users:    auth.NewUserStore(),
		sessions: auth.NewSessionManager(),
		password: password,
	}

	// BCRYPT: If a password is provided, register the default admin user.
	// The password is hashed with bcrypt before storage — plaintext is never kept.
	if password != "" {
		if err := s.users.RegisterUser("admin", password); err != nil {
			fmt.Printf("Warning: failed to register admin user: %v\n", err)
		}
	}

	return s
}

// loadLogs reads and parses all log files.
func (s *Server) loadLogs() error {
	mergedLines, fileErrs := source.ReadMultipleConcurrent(s.files)
	p := parser.SimpleParser{}

	for line := range mergedLines {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		event, err := p.ParseLogLine(parts[1], parts[0])
		if err != nil {
			continue
		}
		s.storage.AddEvent(event)
		s.events = append(s.events, event)
	}

	for err := range fileErrs {
		if err != nil {
			fmt.Printf("Source error: %v\n", err)
		}
	}

	ordering.SortEvents(s.events)
	return nil
}

// ============================================================
// BCRYPT: Authentication middleware
// Checks for a valid session cookie before allowing access.
// If no password is set, all requests are allowed through.
// ============================================================
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// No password set — skip authentication
		if s.password == "" {
			next(w, r)
			return
		}

		// Check for session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// BCRYPT/SESSION: Validate the session token
		_, valid := s.sessions.ValidateSession(cookie.Value)
		if !valid {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}

// ============================================================
// BCRYPT: Login handler — POST /api/login
// Receives username + password, verifies with bcrypt, creates session.
// ============================================================
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Serve login page
		http.ServeFile(w, r, "web/static/login.html")
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	// BCRYPT: Authenticate — compares plaintext against stored bcrypt hash
	if !s.users.Authenticate(username, password) {
		// Wrong credentials — redirect back with error
		http.Redirect(w, r, "/login?error=1", http.StatusFound)
		return
	}

	// BCRYPT: Authentication successful — create a session token
	token, err := s.sessions.CreateSession(username)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // prevent XSS access to cookie
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

// handleLogout destroys the session and redirects to login.
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		s.sessions.DestroySession(cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

// handleLogs handles GET /api/logs with optional ?level= and ?match= query params.
func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	level := r.URL.Query().Get("level")
	match := r.URL.Query().Get("match")

	var filters []filter.Filter

	if level != "" {
		minLevel := models.ParseLogLevel(level)
		if minLevel != models.UNKNOWN {
			filters = append(filters, filter.LevelFilter{MinLevel: minLevel})
		}
	}

	if match != "" {
		filters = append(filters, filter.KeywordFilter{Keyword: match})
	}

	result := s.events
	if len(filters) > 0 {
		activeFilter := filter.ChainFilters(filters...)
		result = filter.ApplyFilter(result, activeFilter)
	}

	resp := APIResponse{
		Total:  len(result),
		Events: result,
	}

	json.NewEncoder(w).Encode(resp)
}

// handleStats handles GET /api/stats.
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stats := &StatsResponse{
		TotalEvents: s.storage.TotalEvents(),
		LevelCounts: s.storage.GetLevelStatistics(),
		Sources:     s.storage.GetAllSources(),
	}

	resp := APIResponse{
		Total: stats.TotalEvents,
		Stats: stats,
	}

	json.NewEncoder(w).Encode(resp)
}

// Start loads logs and starts the HTTP server.
func (s *Server) Start(port int) error {
	if err := s.loadLogs(); err != nil {
		return err
	}

	mux := http.NewServeMux()

	// Auth routes (always accessible)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/logout", s.handleLogout)

	// Protected API routes — wrapped with auth middleware
	mux.HandleFunc("/api/logs", s.requireAuth(s.handleLogs))
	mux.HandleFunc("/api/stats", s.requireAuth(s.handleStats))

	// Protected static files (dashboard)
	fs := http.FileServer(http.Dir("web/static"))
	mux.HandleFunc("/", s.requireAuth(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

	addr := fmt.Sprintf(":%d", port)
	if s.password != "" {
		fmt.Printf("\n🔐 Dashboard running at http://localhost%s (password protected)\n", addr)
		fmt.Printf("   Login with username: admin\n\n")
	} else {
		fmt.Printf("\n🌐 Dashboard running at http://localhost%s\n\n", addr)
	}

	return http.ListenAndServe(addr, mux)
}
