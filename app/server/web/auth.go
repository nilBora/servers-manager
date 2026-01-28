package web

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nilBora/servers-manager/app/store"
)

const (
	sessionCookieName = "session_id"
	sessionDuration   = 7 * 24 * time.Hour // 7 days
	bcryptCost        = 12
)

type contextKey string

const userContextKey contextKey = "user"

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSessionID creates a random session ID
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AuthMiddleware checks if user is authenticated
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session cookie
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate session
		session, err := h.store.GetSession(r.Context(), cookie.Value)
		if err != nil {
			// Clear invalid cookie
			http.SetCookie(w, &http.Cookie{
				Name:     sessionCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get user
		user, err := h.store.GetUserByID(r.Context(), session.UserID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetCurrentUser returns the current user from context
func GetCurrentUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		return nil
	}
	return user
}

// handleLogin renders the login page
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		if _, err := h.store.GetSession(r.Context(), cookie.Value); err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	// Check if setup is needed (no users exist)
	count, _ := h.store.CountUsers(r.Context())
	if count == 0 {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}

	data := templateData{
		Theme: h.getTheme(r),
	}

	if err := h.tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleLoginPost handles login form submission
func (h *Handler) handleLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderLoginError(w, r, "Invalid form data")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		h.renderLoginError(w, r, "Username and password are required")
		return
	}

	// Get user
	user, err := h.store.GetUserByUsername(r.Context(), username)
	if err != nil {
		h.renderLoginError(w, r, "Invalid username or password")
		return
	}

	// Check password
	if !CheckPassword(password, user.PasswordHash) {
		h.renderLoginError(w, r, "Invalid username or password")
		return
	}

	// Create session
	sessionID, err := GenerateSessionID()
	if err != nil {
		h.renderLoginError(w, r, "Failed to create session")
		return
	}

	session := &store.Session{
		ID:        sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(sessionDuration),
	}

	if err := h.store.CreateSession(r.Context(), session); err != nil {
		h.renderLoginError(w, r, "Failed to create session")
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(sessionDuration.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to dashboard
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handleLogout logs out the user
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Get and delete session
	if cookie, err := r.Cookie(sessionCookieName); err == nil {
		_ = h.store.DeleteSession(r.Context(), cookie.Value)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// handleSetup renders the initial setup page
func (h *Handler) handleSetup(w http.ResponseWriter, r *http.Request) {
	// Check if setup already done
	count, _ := h.store.CountUsers(r.Context())
	if count > 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := templateData{
		Theme: h.getTheme(r),
	}

	if err := h.tmpl.ExecuteTemplate(w, "setup.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleSetupPost handles initial setup form submission
func (h *Handler) handleSetupPost(w http.ResponseWriter, r *http.Request) {
	// Check if setup already done
	count, _ := h.store.CountUsers(r.Context())
	if count > 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.renderSetupError(w, r, "Invalid form data")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if username == "" || password == "" {
		h.renderSetupError(w, r, "Username and password are required")
		return
	}

	if len(password) < 6 {
		h.renderSetupError(w, r, "Password must be at least 6 characters")
		return
	}

	if password != confirmPassword {
		h.renderSetupError(w, r, "Passwords do not match")
		return
	}

	// Hash password
	hash, err := HashPassword(password)
	if err != nil {
		h.renderSetupError(w, r, "Failed to create user")
		return
	}

	// Create user
	user := &store.User{
		Username:     username,
		PasswordHash: hash,
	}

	if err := h.store.CreateUser(r.Context(), user); err != nil {
		if errors.Is(err, store.ErrConflict) {
			h.renderSetupError(w, r, "Username already exists")
			return
		}
		h.renderSetupError(w, r, "Failed to create user")
		return
	}

	// Redirect to login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) renderLoginError(w http.ResponseWriter, r *http.Request, message string) {
	data := templateData{
		Theme: h.getTheme(r),
		Error: message,
	}
	w.WriteHeader(http.StatusBadRequest)
	_ = h.tmpl.ExecuteTemplate(w, "login.html", data)
}

func (h *Handler) renderSetupError(w http.ResponseWriter, r *http.Request, message string) {
	data := templateData{
		Theme: h.getTheme(r),
		Error: message,
	}
	w.WriteHeader(http.StatusBadRequest)
	_ = h.tmpl.ExecuteTemplate(w, "setup.html", data)
}
