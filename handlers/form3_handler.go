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
		pathParams := mux.Vars(r)
		accountId := pathParams["accountId"]
		if accountId == "" {
			http.Error(w, errors.Wrap(nil, "Missing 'accountId' param").Error(), http.StatusBadRequest)
			return
		}
		if account, err = form3Client.GetAccount(accountId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(account); err != nil {
			http.Error(w, errors.Wrap(err, "Could not encode account into json").Error(), http.StatusInternalServerError)
		}
		return
	}
}

func DeleteAccount(form3Client form3_client.Form3ClientIface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
		)
		log.Println("inside delete")
		params := mux.Vars(r)
		accountId := params["accountId"]
		if accountId == "" {
			http.Error(w, errors.Wrap(nil, "Missing 'accountId' param").Error(), http.StatusBadRequest)
			return
		}
		version := r.URL.Query().Get("version")
		log.Println(version)
		if version == "" {
			http.Error(w, errors.Wrap(nil, "Missing 'version' param").Error(), http.StatusBadRequest)
			return
		}

		if err = form3Client.DeleteAccount(accountId, version); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusNoContent)
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
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(account); err != nil {
			http.Error(w, errors.Wrap(err, "Could not encode account into json").Error(), http.StatusInternalServerError)
		}
		return
	}
}
