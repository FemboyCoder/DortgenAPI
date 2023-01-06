package main

import (
	"DortgenAPI/src/api"
	"DortgenAPI/src/database"
	"DortgenAPI/src/util"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	APIPort          = flag.String("port", "3000", "port to host the api on")
	GenerateCooldown = flag.Int("cooldown", 300, "cooldown in seconds for generating alts")
	router           chi.Router
)

func init() {
	router = chi.NewRouter()
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			if request.URL.Path == "/favicon.ico" {
				handler.ServeHTTP(writer, request)
				return
			}

			f := &middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: false}
			entry := f.NewLogEntry(request)
			ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)
			t1 := time.Now()
			defer func() {
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()

			handler.ServeHTTP(ww, request)
		})
	})

	err := setupEndpoints()
	if err != nil {
		log.Fatal("Error setting up API endpoints: " + err.Error())
	}

}

func main() {

	datapath := "dortgenapi"

	var err error

	// create the data folder if it doesn't exist
	err = util.CreateFolderIfNotExists(datapath)
	if err != nil {
		log.Fatal("Error creating data folder: " + err.Error())
	}

	// starts the database connection and sets up the tables
	err = database.Startup(datapath, *GenerateCooldown)
	if err != nil {
		log.Fatal("Error starting up database: " + err.Error())
	}

	// test creating an api key
	key, err := database.CreateApiKey("royalty")
	if err != nil {
		log.Println("Error creating api key: " + err.Error())
	}
	log.Println(*key)

	// start the web api
	log.Println("Listening for API requests at port " + *APIPort)
	err = http.ListenAndServe(":"+*APIPort, router)
	if err != nil {
		log.Fatal(err)
	}

}

func setupEndpoints() error {

	router.Get("/favicon.ico", func(writer http.ResponseWriter, request *http.Request) {
		// returns the favicon.ico file
		http.ServeFile(writer, request, "public/favicon.ico")
	})

	router.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		// get the content type from request
		contentType := request.Header.Get("Content-Type")
		// if the content type is application/json, return a json response
		if contentType == "application/json" {
			writer.WriteHeader(http.StatusNotFound)
			_, err := writer.Write([]byte(`{"error": "not found"}`))
			if err != nil {
				log.Println("Error writing response: " + err.Error())
			}
			return
		}
		// otherwise, return a html response
		writer.WriteHeader(http.StatusNotFound)
		http.ServeFile(writer, request, "public/404.html")
	})

	router.Get(
		"/status",
		api.StatusFunc,
	)

	router.Get(
		"/generate",
		api.GenerateFunc,
	)

	router.Get("/validate", api.ValidateFunc)

	return nil
}
