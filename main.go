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
	log.Println("in handle")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/form3Client/account/{accountId}", handlers.GetAccount(form3Client))
	myRouter.HandleFunc("/form3Client/account", handlers.CreateAccount(form3Client))
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
