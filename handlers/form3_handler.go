package handlers

import (
	"encoding/json"
	form3_client "form3-interview/clients"
	"form3-interview/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func GetAccount(form3Client form3_client.Form3ClientIface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			account models.AccountWrapper
			err     error
		)
		w.Header().Set("Content-Type", "application/json")
		pathParams := mux.Vars(r)
		accountId, ok := pathParams["accountId"]
		if !ok {
			http.Error(w, errors.Wrap(errors.New("validation"), "Missing 'accountId' param").Error(), http.StatusBadRequest)
			return
		}
		if account, err = form3Client.GetAccount(accountId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(account); err != nil {
			http.Error(w, errors.Wrap(err, "Could not encode account into json").Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func DeleteAccount(form3Client form3_client.Form3ClientIface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
		)

		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		accountId, ok := params["accountId"]
		if !ok {
			http.Error(w, errors.Wrap(errors.New("validation"), "Missing 'accountId' param").Error(), http.StatusBadRequest)
			return
		}

		version := r.URL.Query().Get("version")
		if len(version) == 0 {
			http.Error(w, errors.Wrap(errors.New("validation"), "Missing 'version' param").Error(), http.StatusBadRequest)
			return
		}
		if err = form3Client.DeleteAccount(accountId, version); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

		w.Header().Set("Content-Type", "application/json")
		if account, err = form3Client.PostAccount(r.Body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = json.NewEncoder(w).Encode(account); err != nil {
			http.Error(w, errors.Wrap(err, "Could not encode account into json").Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
}
