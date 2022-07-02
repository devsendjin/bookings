package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/devsendjin/bookings/pkg/config"
	"github.com/devsendjin/bookings/pkg/handlers"
	"github.com/devsendjin/bookings/pkg/render"
)

const appPort = ":8081"

var app config.AppConfig
var session *scs.SessionManager

// main in the main application function
func main() {
	app.IsDevelopment = true
	app.IsProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache
	app.UseCache = app.IsProduction

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("\nStarting application on port %s\n\n\n", appPort)

	server := &http.Server{
		Addr:    appPort,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}
