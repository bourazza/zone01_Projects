package handlers

import (
	"groupie/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var templates *template.Template

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"ErrorCode":    "400",
		"ErrorMessage": "Bad Request",
	})
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := models.FetchArtists()
	if err != nil {
		ServerErrorHandler(w, r)
		return
	}
	templates.ExecuteTemplate(w, "index.html", artists)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"ErrorCode":    "404",
		"ErrorMessage": "Page Not Found",
	})
}

func ServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"ErrorCode":    "500",
		"ErrorMessage": "Internal Server Error",
	})
}
func ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"ErrorCode":    "403",
		"ErrorMessage": "Forbidden Access",
	})
}

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		BadRequestHandler(w, r)
		return
	}

	artistID, err := strconv.Atoi(id)
	if err != nil {
		BadRequestHandler(w, r)
		return
	}

	selectedArtist, err := models.FetchArtistByID(artistID)
	if err != nil {
		ServerErrorHandler(w, r)
		return
	}

	if selectedArtist == nil {
		NotFoundHandler(w, r)
		return
	}

	log.Printf("Artist %s has locations: %+v", selectedArtist.Name, selectedArtist.Locations)
	log.Printf("Artist %s has concert dates: %+v", selectedArtist.Name, selectedArtist.ConcertDates)

	templates.ExecuteTemplate(w, "artist.html", map[string]interface{}{
		"Artist": selectedArtist,
		"ID":     artistID,
	})
}


// UnauthorizedHandler handles unauthorized access attempts
func UnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"ErrorCode":    "401",
		"ErrorMessage": "Unauthorized Access",
	})
}