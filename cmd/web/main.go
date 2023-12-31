package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/atrifunac7/bookings/pkg/config"
	"github.com/atrifunac7/bookings/pkg/handlers"
	"github.com/atrifunac7/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var session *scs.SessionManager
var app config.AppConfig

// main is the main function of the application
func main() {
	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
