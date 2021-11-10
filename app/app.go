package app

import (
	form3_client "form3-interview/clients"
	"form3-interview/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	Client form3_client.Form3Client
}

var app *App

func NewApp() *App {
	return &App{
		Router: mux.NewRouter().StrictSlash(true),
		Client: form3_client.Form3Client{},
	}
}

func (a *App) HandleRequests() {
	if app == nil {
		app = NewApp()
	}
	log.Println("inside app")
	app.Router.HandleFunc("/form3Client/accounts/{accountId}", handlers.GetAccount(app.Client)).Methods(http.MethodGet)
	app.Router.HandleFunc("/form3Client/accounts", handlers.CreateAccount(app.Client)).Methods(http.MethodPost)
	app.Router.HandleFunc("/form3Client/accounts/{accountId}", handlers.DeleteAccount(app.Client)).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":10000", app.Router))
}
