/*
In this code, we use the http.ListenAndServeTLS function to start the server with HTTPS.
We pass the path to the SSL/TLS certificate and key as arguments to this function.

To create an SSL/TLS certificate and key, you can use tools like OpenSSL.
Here's an example of how you can create a self-signed certificate and key:

openssl req -x509 -newkey rsa:2048 -nodes -keyout key.pem -out cert.pem -days 365 -subj "/C=US/ST=State/L=Locality/O=Organization/CN=localhost"

This will create a self-signed certificate and key in the cert.pem and key.pem files, respectively.

Note that self-signed certificates are not trusted by default by most browsers
and clients, so you may need to add an exception or install the certificate as a trusted root certificate.

Alternatively, you can obtain a trusted certificate from a certificate authority (CA) like Let's Encrypt.

To use Let's Encrypt, you can use the certbot tool to obtain a certificate and key. Here's an example
of how you can use certbot to obtain a certificate and key:

certbot certonly --standalone -d example.com

This will create a certificate and key in the /etc/letsencrypt/live/example.com directory.

You can then use the cert.pem and privkey.pem files from this directory as your SSL/TLS certificate and key.

Make sure to replace the path/to/cert.pem and path/to/key.pem placeholders with the actual paths to your SSL/TLS certificate and key files.

*/

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"gorm.io/gorm"

	config "epg/src/config"
	"epg/src/controller"

	"github.com/gorilla/mux"
)

type PageData struct {
	Title   string
	Heading string
	Content string
}

type Server struct {
	mux  *mux.Router
	port int
	db   *gorm.DB
	srv  *http.Server
}

func NewServer(port int, db *gorm.DB) *Server {
	return &Server{
		mux:  mux.NewRouter(),
		port: port,
		db:   db,
	}
}

func (s *Server) setupRoutes() {

	// HTML Routes index
	s.mux.HandleFunc("/", indexHandler)
	s.mux.HandleFunc("/contact", contactHandler)

	// Network routes
	networkHandler := controller.NewNetworkHandler(s.db)
	s.mux.HandleFunc("/network", networkHandler.GetAllNetworksHTML).Methods("GET")
	s.mux.HandleFunc("/network/{networkId}", networkHandler.GetNetworkByIdHTML).Methods("GET")

	// Channels routes
	channelHandler := controller.NewChannelHandler(s.db)
	s.mux.HandleFunc("/channel", channelHandler.GetAllChannelsHTML).Methods("GET")
	s.mux.HandleFunc("/channel/{channelId}", channelHandler.GetChannelByIdHTML).Methods("GET")

	// Genre routes
	genreHandler := controller.NewGenreHandler(s.db)
	s.mux.HandleFunc("/genre", genreHandler.GetAllGenresHTML).Methods("GET")
	s.mux.HandleFunc("/genre/{genreId}", genreHandler.GetGenreByIdHTML).Methods("GET")

	// Country routes
	countryHandler := controller.NewCountryHandler(s.db)
	s.mux.HandleFunc("/country", countryHandler.GetAllCountries).Methods("GET")
	s.mux.HandleFunc("/country/{countryId}", countryHandler.GetCountryById).Methods("GET")
	s.mux.HandleFunc("/countrycode/{countryCode}", countryHandler.GetCountryByCode).Methods("GET")
	s.mux.HandleFunc("/country", countryHandler.CreateCountry).Methods("POST")
	s.mux.HandleFunc("/country/{countryId}", countryHandler.UpdateCountry).Methods("PUT")
	s.mux.HandleFunc("/country/{countryId}", countryHandler.DeleteCountry).Methods("DELETE")

	// Category routes
	categoryHandler := controller.NewCategoryHandler(s.db)
	s.mux.HandleFunc("/category", categoryHandler.GetAllCategoriesHTML).Methods("GET")
	s.mux.HandleFunc("/category/{categoryId}", categoryHandler.GetCategoryByIdHTML).Methods("GET")

	// Timezone routes
	timezoneHandler := controller.NewTimezoneHandler(s.db)
	s.mux.HandleFunc("/timezone", timezoneHandler.GetAllTimezones).Methods("GET")
	s.mux.HandleFunc("/timezone/{timezoneId}", timezoneHandler.GetTimezoneById).Methods("GET")

	// Rating system routes
	ratingHandler := controller.NewRatingHandler(s.db)
	s.mux.HandleFunc("/ratingsystem", ratingHandler.GetAllRatingSystems).Methods("GET")
	s.mux.HandleFunc("/rating", ratingHandler.GetAllRatingsHTML).Methods("GET")
	s.mux.HandleFunc("/rating/{ratingId}", ratingHandler.GetRatingSystemById).Methods("GET")

	// Rating value routes
	ratingValueHandler := controller.NewRatingValueHandler(s.db)
	s.mux.HandleFunc("/ratingvalue", ratingValueHandler.GetAllRatingValues).Methods("GET")
	s.mux.HandleFunc("/ratingvalue/{ratingValueId}", ratingValueHandler.GetRatingValueById).Methods("GET")

	// Events routes
	eventHandler := controller.NewEventHandler(s.db)
	s.mux.HandleFunc("/event", eventHandler.GetAllEvents).Methods("GET")
	s.mux.HandleFunc("/event/{eventId}", eventHandler.GetEventById).Methods("GET")
	s.mux.HandleFunc("/event", eventHandler.CreateEvent).Methods("POST")
	s.mux.HandleFunc("/event/{eventId}", eventHandler.UpdateEvent).Methods("PUT")
	s.mux.HandleFunc("/event/{eventId}", eventHandler.DeleteEvent).Methods("DELETE")

	// Event rating routes
	eventRatingHandler := controller.NewEventRatingHandler(s.db)
	s.mux.HandleFunc("/eventrating", eventRatingHandler.GetAllEventRatings).Methods("GET")
	s.mux.HandleFunc("/eventrating/{eventId}/{ratingValueId}", eventRatingHandler.GetEventRatingById).Methods("GET")
	s.mux.HandleFunc("/eventbytime/{time}", eventHandler.GetEventByTime).Methods("GET")
	s.mux.HandleFunc("/eventrating", eventRatingHandler.CreateEventRating).Methods("POST")
	s.mux.HandleFunc("/eventrating/{eventId}/{ratingValueId}", eventRatingHandler.UpdateEventRating).Methods("PUT")
	s.mux.HandleFunc("/eventrating/{eventId}/{ratingValueId}", eventRatingHandler.DeleteEventRating).Methods("DELETE")

	// Serve static files
	staticDir := http.Dir("./static")
	s.mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Serve JavaScript files
	jsDir := http.Dir("./static/js")
	s.mux.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(jsDir)))

	// Serve CSS files
	cssDir := http.Dir("./static/css")
	s.mux.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(cssDir)))

	// Serve images
	imagesDir := http.Dir("./static/img")
	s.mux.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(imagesDir)))

	// Serve HTML templates
	htmlDir := http.Dir("./static/html")
	s.mux.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(htmlDir)))

	// Serve SSL/TLS certificates
	//certDir := http.Dir("./static/cert")
	//s.mux.PathPrefix("/cert/").Handler(http.StripPrefix("/cert/", http.FileServer(certDir)))

	// Serve favicon.ico
	s.mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/icon/favicon.ico")
	})
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Index Page",
		Heading: "Welcome to EPG",
		Content: "This is the index page content.",
	}

	controller.RenderTemplate(w, "static/html/index.html", data)
}

// ContactHandler handles the contact form submission
// contactHandler handles the contact form submission
func contactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		// Send an email to the specified recipient
		err := sendEmail("recipient-email@example.com", "Contact Form Submission", "Name: "+name+"\nEmail: "+email+"\nMessage: "+message)
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		// Render a success message
		data := PageData{
			Title:   "Contact",
			Heading: "Contact Us",
			Content: "Thank you for contacting us! We will respond to your message as soon as possible.",
		}
		controller.RenderTemplate(w, "static/html/contact.html", data)
	} else {
		// Render the contact form
		data := PageData{
			Title:   "Contact",
			Heading: "Contact Us",
		}
		controller.RenderTemplate(w, "static/html/contact.html", data)
	}
}

func (s *Server) start() error {
	// Create server with timeouts
	s.srv = &http.Server{
		Handler:      s.mux,
		Addr:         fmt.Sprintf(":%d", s.port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	// Start server with HTTPS
	//return s.srv.ListenAndServeTLS("./cert/cert.pem", "./cert/key.pem")
	return s.srv.ListenAndServe()
}

func main() {
	// Connect to the database

	db := controller.GetDB()

	// Create a new server
	server := NewServer(config.Config.BindPort, db)

	// Setup routes
	server.setupRoutes()
	/*
		go func() {
			for {
				// Wait until 0:00 the next day
				now := time.Now()
				nextDay := now.Add(24 * time.Hour)
				nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
				time.Sleep(nextDay.Sub(now))

				// Populate the database with initial events
				err := (&Event{}).PopulateInitialEvents(dbConn, []uint{1, 2}, nextDay)
				if err != nil {
					log.Println(err)
				}
			}
		}()
	*/
	// Start the server
	log.Println("Listening: port = " + strconv.Itoa(server.port))
	log.Fatal(server.start())
	log.Println("Shutting down...")
}

// Send an email using a mail server
func sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", "your-email@gmail.com", "your-password", "smtp.gmail.com")
	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587", auth, "your-email@gmail.com", []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
