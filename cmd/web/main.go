package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/trcwrx78/bookings/internal/config"
	"github.com/trcwrx78/bookings/internal/driver"
	"github.com/trcwrx78/bookings/internal/handlers"
	"github.com/trcwrx78/bookings/internal/helpers"
	"github.com/trcwrx78/bookings/internal/models"
	"github.com/trcwrx78/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// The main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server {
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
		// What am I going to put in the session
		gob.Register(models.Reservation{})
		gob.Register(models.User{})
		gob.Register(models.Room{})
		gob.Register(models.Restriction{})

		// change this to true in production
		app.InProduction = false

		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		app.InfoLog = infoLog

		errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		app.ErrorLog = errorLog

		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
	
		app.Session = session

		// connect to database
		log.Println("Connecting to database...")
		db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=BrynAd02")
		if err != nil {
			log.Fatal("Cannot connect to database! Dying...")
		}
		log.Println("Connected to the database!")
	
		tc, err := render.CreateTemplateCache()
		if err !=nil {
			log.Fatal("cannot create template cache")
			return nil, err
		}
	
		app.TemplateCache = tc
		app.UseCache = false
	
		repo := handlers.NewRepo(&app, db)
		handlers.NewHandlers(repo)
		render.NewRenderer(&app)
		helpers.NewHelpers(&app)
	

	return db, nil
}