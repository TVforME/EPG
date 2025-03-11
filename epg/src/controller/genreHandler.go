package controller

import (
	"net/http"
	"strconv"

	"epg/src/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GenreHandler manages genre endpoints.
type GenreHandler struct {
	db *gorm.DB
}

// NewGenreHandler returns a new instance of GenreHandler.
func NewGenreHandler(db *gorm.DB) *GenreHandler {
	return &GenreHandler{db: db}
}

// GetAllGenresHTML fetches genres from the /genre endpoint and renders the HTML page.
func (gh *GenreHandler) GetAllGenresHTML(w http.ResponseWriter, r *http.Request) {
	var genres []model.Genre

	// Retrieve all genres with a left join to the genrecolor table.
	if err := gh.db.Preload("GenreColor").Find(&genres).Error; err != nil {
		HandleHtmlError(w, err)
		return
	}
	// TODO work out why 0 returns empty string
	// Build response by iterating over genres.
	var responses []model.Genre
	for _, g := range genres {
		// Set the ColorHex to the default color if it's empty.
		if g.GenreColor == nil || g.GenreColor.ColorHex == "" {
			g.ColorHex = "#A0A0A0"
		} else {
			g.ColorHex = g.GenreColor.ColorHex
		}

		responses = append(responses, g)
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:   "Genre",
		Heading: "Genre List",
		Genres:  responses,
	}

	RenderTemplate(w, "static/html/genre.html", data)
}

// GetGenreById handles GET requests for retrieving a genre by its ID.
func (gh *GenreHandler) GetGenreByIdHTML(w http.ResponseWriter, r *http.Request) {

	genreID, err := strconv.Atoi(mux.Vars(r)["genreId"])
	if err != nil {
		HandleHtmlError(w, err)
		return
	}

	var genre model.Genre
	// Retrieve genre with preloading GenreColor.
	if err = gh.db.Preload("GenreColor").Where("genre_id = ?", genreID).First(&genre).Error; err != nil {
		HandleHtmlError(w, err)
		return
	}
	// TODO work out why 0 returns empty string
	// Set the ColorHex to a default value if it's empty.
	if genre.GenreColor == nil || genre.GenreColor.ColorHex == "" {
		genre.ColorHex = "#606060"
	} else {
		genre.ColorHex = genre.GenreColor.ColorHex
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:   "Genre",
		Heading: "Genre ByID",
		Genres:  []model.Genre{genre},
	}

	RenderTemplate(w, "static/html/genre.html", data)
}
