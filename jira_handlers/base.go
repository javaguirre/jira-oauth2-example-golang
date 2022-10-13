package jira_handlers

import (
	"net/http"
)

func New() http.Handler {
	mux := http.NewServeMux()
	// Root
	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	// OAuth
	mux.HandleFunc("/auth/jira/login", oauthJiraLogin)
	mux.HandleFunc("/auth/jira/callback", oauthJiraCallback)

	return mux
}
