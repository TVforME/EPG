// categoryHandler.go
package controller

import (
	"net/http"
	"strconv"

	"epg/src/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CatergoryResponse is the structure returned via API responses.
type CategoryResponse struct {
	CategoryID  uint   `json:"CategoryID"`
	Description string `json:"description"`
}

// CategoryHandler ...
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler ...
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// GetAllCategoriesHTML fetches category from the /catergory endpoint and renders the HTML page.
func (ch *CategoryHandler) GetAllCategoriesHTML(w http.ResponseWriter, r *http.Request) {

	categories := []model.Category{}

	err := ch.db.Preload("Events").Find(&categories).Error
	if err != nil {
		HandleHtmlError(w, err)
		return
	}

	var responses []model.Category
	for _, c := range categories {
		responses = append(responses, model.Category{
			CategoryID:  c.CategoryID,
			Description: c.Description,
		})

	}

	// Prepare data for template rendering.
	data := PageData{
		Title:    "Category",
		Heading:  "Category List",
		Category: responses,
	}

	RenderTemplate(w, "static/html/category.html", data)
}

// GetCategoryByIdHTML handles GET requests for retrieving a category by its ID.
func (ch *CategoryHandler) GetCategoryByIdHTML(w http.ResponseWriter, r *http.Request) {

	categoryId, err := strconv.Atoi(mux.Vars(r)["categoryId"])
	if err != nil {
		HandleHtmlError(w, err)
		return
	}

	var category model.Category
	err = ch.db.Preload("Events").Where("category_id = ?", categoryId).First(&category).Error
	if err != nil {
		HandleHtmlError(w, err)
		return
	}

	response := model.Category{
		CategoryID:  category.CategoryID,
		Description: category.Description,
	}

	// Prepare data for template rendering.
	data := PageData{
		Title:    "Category",
		Heading:  "Category ByID",
		Category: []model.Category{response},
	}

	RenderTemplate(w, "static/html/category.html", data)
}
