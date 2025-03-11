# EPG
EPG system written in Go using SQLite and REST backend.

This is not a completed project and partically implemented with the required events table to be laoded and queried.

below is the file structure:

epg
|   build-and-run.sh
|   go.mod
|   go.sum
|   README.md
|   details.txt
|   
+---bin
|       epg.db
|       
+---cert
+---csv
|       categories.csv
|       color.csv
|       countries.csv
|       events_template.csv
|       genre.csv
|       ratings.csv
|       ratingsystems.csv
|       timezones.csv
|       
+---notes
|       CRID info.doc
|       
+---src
|   +---app
|   |       main.go
|   |       
|   +---config
|   |       config.go
|   |       config.json
|   |       
|   +---controller
|   |       base.go
|   |       categoryHandler.go
|   |       channelHandler.go
|   |       countryHandler.go
|   |       eventHandler.go
|   |       eventRatingHandler.go
|   |       genreHandler.go
|   |       networkHandler.go
|   |       ratingHandler.go
|   |       ratingValueHandler.go
|   |       timezoneHandler.go
|   |       
|   \---model
|           categoryModel.go
|           channelModel.go
|           countryModel.go
|           eventModel.go
|           genreModel.go
|           model.go
|           networkModel.go
|           ratingsModel.go
|           timezoneModel.go
|           
\---static
    +---css
    |       styles.css
    |       
    +---html
    |       about.html
    |       base.html
    |       category.html
    |       channel.html
    |       contact.html
    |       epg.html
    |       genre.html
    |       index.html
    |       network.html
    |       rating.html
    |       
    +---icon
    |       favicon.ico
    |       
    +---img
    |   +---logos
    |   |       default-50x50.png
    |   |       vk3rgl-hd1-50x50.png
    |   |       vk3rgl-hd2-50x50.png
    |   |       
    |   \---ratings
    |       \---au
    |               AV15+.png
    |               C.png
    |               Check.png
    |               E.png
    |               G.png
    |               M.png
    |               MA15+.png
    |               P.png
    |               PG.png
    |               R18+.png
    |               X18+.png
    |               
    \---js
            scripts.js




The project requires a number of Go packages
"gorm.io/gorm"
"github.com/gorilla/mux"


