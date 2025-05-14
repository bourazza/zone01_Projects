package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"groupie/handlers"
)

func main() {
	// Initialize templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Error parsing templates:", err)
	}
	handlers.InitTemplates(templates)

	staticHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// urlPath := strings.TrimPrefix(r.URL.Path, "/static/")
		// filePath := filepath.Join("/", urlPath)

		// if strings.Contains(urlPath, "..") {
		// 	handlers.BadRequestHandler(w, r)
		// 	return
		// }

		filePath := r.URL.Path

		fileInfo, err := os.Stat(filePath[1:])
		if err != nil {
			handlers.NotFoundHandler(w, r)
			return
		}

		if fileInfo.IsDir() {
			handlers.ForbiddenHandler(w, r)
			return
		}

		allowedExtensions := []string{".css", ".js", ".jpg", ".jpeg", ".png", ".gif", ".svg", ".ico"}
		extension := strings.ToLower(filepath.Ext(filePath))

		extensionAllowed := false
		for _, allowed := range allowedExtensions {
			if extension == allowed {
				extensionAllowed = true
				break
			}
		}

		if !extensionAllowed {
			handlers.ForbiddenHandler(w, r)
			return
		}

		if strings.Contains(filePath, "/protected/") {
			handlers.UnauthorizedHandler(w, r)
			return
		}

		http.ServeFile(w, r, filePath[1:])
	})

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handlers.NotFoundHandler(w, r)
			return
		}
		handlers.HomeHandler(w, r)
	})

	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/404", handlers.NotFoundHandler)
	http.HandleFunc("/500", handlers.ServerErrorHandler)
	http.HandleFunc("/400", handlers.BadRequestHandler)
	http.HandleFunc("/401", handlers.UnauthorizedHandler)
	http.HandleFunc("/403", handlers.ForbiddenHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)

	http.HandleFunc("/", rootHandler)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
