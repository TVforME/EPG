/*
The init() function in the controller package is called automatically when the package is initialized.
and establishes the database connection and runs the migration.
When you call controller.GetDB() in your main function, it returns the already established database connection.

*/

package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	config "epg/src/config"
	"epg/src/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

// PageData holds data to render HTML templates.
type PageData struct {
	Title    string
	Heading  string
	Network  []model.Network
	Channels []model.Channel
	Genres   []model.Genre
	Category []model.Category
	Ratings  []model.RatingSystem
}

// Check we have a database in Path.
func checkDatabaseExists(dbPath string) bool {
	_, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal(err)
	}
	return true
}

/*
Each time init() is run we check we have a existing database (assuming nothing has been changed)
If happens the database is deleted, init() creates a new database and Automigrate the models and populate
each table with the relavent values for model csv files.
*/
func init() {
	dbPath := fmt.Sprintf("./bin/%s", config.Config.Dbname)
	var err error
	switch config.Config.DbType {
	case "sqlite3":
		if !checkDatabaseExists(dbPath) {
			log.Println("Database file does not exist. Creating...")
			// Create the database file
			f, err := os.Create(dbPath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			// Open the database and run migration
			dbConn, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
			if err != nil {
				log.Fatal(err)
			}
			// Run migration
			err = dbConn.AutoMigrate(
				&model.Country{},
				&model.Timezone{},
				&model.GenreColor{},
				&model.Genre{},
				&model.Category{},
				&model.RatingSystem{},
				&model.RatingValue{},
				&model.Network{},
				&model.Channel{},
				&model.Event{},
				&model.EventRating{},
			)
			if err != nil {
				log.Fatal(err)
			}

			// Populate the countries table
			country := &model.Country{}
			err = country.LoadFromCSV(dbConn, "countries.csv")
			if err != nil {
				log.Fatal(err)
			}
			// Populate the timezone table
			timezone := &model.Timezone{}
			err = timezone.LoadFromCSV(dbConn, "timezones.csv")
			if err != nil {
				log.Fatal(err)
			}
			// Populate the genre colors table
			genrecolor := &model.GenreColor{}
			err = genrecolor.LoadFromCSV(dbConn, "color.csv")
			if err != nil {
				log.Fatal(err)
			}

			genre := &model.Genre{}
			err = genre.LoadFromCSV(dbConn, "genre.csv")
			if err != nil {
				log.Fatal(err)
			}

			category := &model.Category{}
			err = category.LoadFromCSV(dbConn, "categories.csv")
			if err != nil {
				log.Fatal(err)
			}

			ratingvalue := &model.RatingValue{}
			err = ratingvalue.LoadFromCSV(dbConn, "ratings.csv")
			if err != nil {
				log.Fatal(err)
			}
			ratingsystem := &model.RatingSystem{}
			err = ratingsystem.LoadFromCSV(dbConn, "ratingsystems.csv")
			if err != nil {
				log.Fatal(err)
			}

			// Add our network from config
			err = model.PopulateInitialNetworkValues(dbConn)
			if err != nil {
				log.Fatal(err)
			}

			// Add our channels from config
			err = model.PopulateInitialChannelValues(dbConn)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			// Open the existing database
			dbConn, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
			if err != nil {
				log.Fatal(err)
			}
		}
	default:
		log.Fatal("Unsupported database type")
	}
}

func GetDB() *gorm.DB { return dbConn }

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmplFiles := []string{
		"static/html/base.html",
		tmpl,
	}

	parsedTemplate, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

// handleHtmlError sends an error message with HTTP 500 status.
func HandleHtmlError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
