package main

import (
	form3_client "form3-interview/clients"
	"form3-interview/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {
	form3Client := form3_client.Form3Client{}
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/form3Client/accounts/{accountId}", handlers.GetAccount(form3Client)).Methods(http.MethodGet)
	myRouter.HandleFunc("/form3Client/accounts", handlers.CreateAccount(form3Client)).Methods(http.MethodPost)
	myRouter.HandleFunc("/form3Client/accounts/{accountId}", handlers.DeleteAccount(form3Client)).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
