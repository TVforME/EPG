# EPG
**EPG** (Electronic Program Guide) system written in Go using SQLite and REST backend.

This is not a complete project and is a prototype concept to use with the repeater project.  [DATV Repeater](https://github.com/TVforME/Repeater)
I've since looked into CherryEPG as a subsitute however, I'll continual working on a GO version and support anyone whom has experience and is able to assist with this project.

*DB*
To create a new epg, simply delete the epg.db file in bin folder. On execution a new epg.db file is created using GORM migrate which uses the csv files to populate the tables.
I used this method over static constants to keep the executable smaller in size and to provide flexibility for additions and changes. Some thought on how to package the csv files with the GO install.
I'm yet to decide.

*Ratings*
Each of the ratings systems uses a country identifier (au) here for the rating icon files.

*CRID*
[CRID](https://en.wikipedia.org/wiki/Content_reference_identifier) (content reference identifier) 
I'm yet to work out how modern TV's show the channel icon and current program in their EPG? I'm sure it able to be immplemented.

*TODO*
Implement a query to fill the events table with a full 24 hours of events which is populated once a channel/s is associated to a network. If there is 2 channels to example then need to create 24 hours of events for each channel.
Once events have been created then a query is run to generate a EIT.xml file containing the relavent tags which is injected into the DVB transport stream using TSduck eitinject plugin.
[https://github.com/tsduck/tsduck/tree/master/src/tsplugins](https://github.com/tsduck/tsduck/blob/master/src/tsplugins/tsplugin_eitinject.cpp)


-Program directory folder.

Below is the file structure:
```
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

```


The project requires a number of Go packages
"gorm.io/gorm"
"github.com/gorilla/mux"


