package main

import (
	"fmt"
	"login-form/data"
	"login-form/handlers"
	"net/http"
)

func main() {
	data.InitDB()

	//Home route redirects to login
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/magic-login", handlers.MagicLoginHandler)
	http.HandleFunc("/users", handlers.UsersListHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
