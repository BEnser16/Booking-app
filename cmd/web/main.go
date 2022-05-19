package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BEnser16/go-app/pkg/config"
	"github.com/BEnser16/go-app/pkg/handlers"
	"github.com/BEnser16/go-app/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber string = ":8080"

var session *scs.SessionManager
var app config.AppConfig

func main() {

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

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
