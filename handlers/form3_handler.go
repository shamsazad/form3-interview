package handlers

import (
	"encoding/json"
	form3_client "form3-interview/clients"
	"form3-interview/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	"net/http"
)


func GetAccount(form3Client form3_client.Form3ClientIface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			account models.AccountWrapper
			err     error
		)
		log.Println("inside GetAccount")
		pathParams := mux.Vars(r)
		accountId := pathParams["accountId"]
		log.Println("accountId", accountId)
		if accountId == "" {
			http.Error(w, errors.Wrap(nil, "Missing 'variantId' param").Error(), http.StatusBadRequest)
			return
		}
		if account, err = form3Client.GetAccount(accountId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(account); err != nil {
			http.Error(w, errors.Wrap(err, "Could not encode QbcVariant into json").Error(), http.StatusInternalServerError)
		}
		return
	}
}

func CreateAccount(form3Client form3_client.Form3ClientIface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			account models.AccountWrapper
			err     error
		)

		if account, err = form3Client.PostAccount(r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println(account)
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(account); err != nil {
			http.Error(w, errors.Wrap(err, "Could not encode account into json").Error(), http.StatusInternalServerError)
		}
		return
	}
}

