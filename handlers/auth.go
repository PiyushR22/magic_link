package handlers

import (
	"encoding/json"
	"fmt"
	"login-form/data"
	"login-form/utils"
	"net/http"
	"time"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the ̥̥̥registration HTML page
		http.ServeFile(w, r, "static/register.html")
		return
	}

	if r.Method == http.MethodPost {
		var req RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Username == "" || req.Password == "" {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = data.DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", req.Username, hashedPassword)
		if err != nil {
			http.Error(w, "User already exists or DB error", http.StatusInternalServerError)
			return
		}

		// get the new user ID
		var userID int
		err = data.DB.QueryRow("SELECT id FROM users WHERE username=?", req.Username).Scan(&userID)
		if err != nil {
			http.Error(w, "Failed to fetch user ID", http.StatusInternalServerError)
			return
		}

		// Generate magic login token
		token, err := utils.GenerateRandomToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		expiration := time.Now().AddDate(100, 0, 0)
		_, err = data.DB.Exec("INSERT INTO tokens(user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiration)
		if err != nil {
			http.Error(w, "Failed to save token", http.StatusInternalServerError)
			return
		}

		magicLink := fmt.Sprintf("http://localhost:8080/magic-login?token=%s", token)
		emailBody := fmt.Sprintf("Welcome! Click the link to login: <a href=\"%s\">%s</a>", magicLink, magicLink)

		// Send magic link by email
		err = utils.SendEmail(req.Username, "Your Magic Login Link", emailBody)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "✅ User registered successfully! Check your email for the magic login link.")
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the login HTML form
		http.ServeFile(w, r, "static/login.html")
		return
	}

	if r.Method == http.MethodPost {
		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Username == "" || req.Password == "" {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var userID int
		var hashedPassword string
		err = data.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", req.Username).Scan(&userID, &hashedPassword)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		if !utils.CheckPasswordHash(req.Password, hashedPassword) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateRandomToken()
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		expiration := time.Now().Add(30 * time.Minute)
		_, err = data.DB.Exec("INSERT INTO tokens(user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiration)
		if err != nil {
			http.Error(w, "Failed to save token", http.StatusInternalServerError)
			return
		}

		magicLink := fmt.Sprintf("http://localhost:8080/magic-login?token=%s", token)
		fmt.Println("Magic login link (send this to the user):", magicLink)

		emailBody := fmt.Sprintf("Click the link to login: <a href=\"%s\">%s</a>", magicLink, magicLink)
		err = utils.SendEmail(req.Username, "Your Magic Login Link", emailBody)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "User Login successfully! Check your email for the magic login link. ")

	}

	// Any other HTTP methods
	//	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func MagicLoginHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	var userID int
	var expiresAt time.Time

	err := data.DB.QueryRow("SELECT user_id, expires_at FROM tokens WHERE token = ?", token).Scan(&userID, &expiresAt)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	if time.Now().After(expiresAt) {
		http.Error(w, "Token expired", http.StatusUnauthorized)
		return
	}

	// Optional: Delete token after use so it can't be reused
	//_, _ = data.DB.Exec("DELETE FROM tokens WHERE token = ?", token)

	// User is now "logged in" — you can set session cookie or just respond success
	http.Redirect(w, r, "/users", http.StatusSeeOther)

}

func UsersListHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := data.DB.Query("SELECT username FROM users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	html := `
<!DOCTYPE html>
<html>
<head>
  <title>Registered Users</title>
  <link rel="stylesheet" href="/static/users.css">
</head>
<body>
  <h1>Registered Users</h1>
  <ul>
`

	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			http.Error(w, "Error scanning user", http.StatusInternalServerError)
			return
		}
		html += fmt.Sprintf("<li>%s</li>", username)
	}

	html += `
  </ul>
</body>
</html>
`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
