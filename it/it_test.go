package it_test

import (
	"bytes"
	form3_client "form3-interview/clients"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func Test_form3ClientGet(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		givenPayload interface{}
		err          error
		accountId    string
		status       int
	}{
		{
			name:      "Account id is not uuid",
			accountId: "4567",
			status:    http.StatusBadRequest,
			err:       errors.New("{\"error_message\":\"id is not a valid uuid\"}"),
		},
		{
			name:      "Account doesn't exist",
			accountId: "02d3792a-1c45-4d91-98d0-ca83790afe89",
			status:    http.StatusNotFound,
			err:       errors.New("{\"error_message\":\"record 02d3792a-1c45-4d91-98d0-ca83790afe89 does not exist\"}"),
		},
		{
			name:      "happy path, account retrieved",
			accountId: "cb1e2074-1056-4b27-b4e0-ed9f0c46b067",
			status:    http.StatusOK,
			err:       nil,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := form3_client.Form3Client{
				HttpClient: &http.Client{},
				BaseURL:    getEnv("BASE_URL", "http://localhost:8080/"),
			}

			if test.name == "happy path, account retrieved" {
				createDummyAccount()
			}

			account, err := client.GetAccount(test.accountId)
			if test.err == nil {
				assert.Equal(t, account.Account.ID, test.accountId)
				deleteDummyAccount(test.accountId)
			} else {
				assert.Equal(t, test.err.Error(), err.Error.Error())
				assert.Equal(t, test.status, err.Code)
			}
		})
	}
}

func Test_form3ClientPost(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		givenPayload io.Reader
		err          error
		status       int
		accountId    string
	}{
		{
			name:         "Missing required field, organization_id",
			givenPayload: strings.NewReader("{}"),
			err:          errors.New("Validation error"),
			status:       http.StatusInternalServerError,
		},
		{
			name:         "Bad data sent to server",
			givenPayload: strings.NewReader(""),
			err:          errors.New("Validation error"),
			status:       http.StatusBadRequest,
		},
		{
			name:         "happy path, account created",
			givenPayload: getBody(),
			err:          nil,
			accountId:    "cb1e2074-1056-4b27-b4e0-ed9f0c46b067",
			status:       http.StatusCreated,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := form3_client.Form3Client{
				HttpClient: &http.Client{},
				BaseURL:    getEnv("BASE_URL", "http://localhost:8080/"),
			}

			account, err := client.PostAccount(test.givenPayload)
			if test.err == nil {
				assert.Equal(t, account.Account.ID, test.accountId)
				deleteDummyAccount(test.accountId)
			} else {
				assert.Equal(t, test.err.Error(), err.Message)
				assert.Equal(t, test.status, err.Code)
			}
		})
	}
}

func Test_form3ClientDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		accountId string
		version   string
		err       error
		code      int
	}{
		{
			name:      "Missing required field, account_id and version",
			accountId: "",
			version:   "",
			err:       errors.New("{\"code\":\"PAGE_NOT_FOUND\",\"message\":\"Page not found\"}"),
			code:      http.StatusNotFound,
		},
		{
			name:      "Account id not uuid",
			accountId: "1234",
			version:   "0",
			err:       errors.New("{\"error_message\":\"id is not a valid uuid\"}"),
			code:      http.StatusBadRequest,
		},
		{
			name:      "Missing required field, version",
			accountId: "cb1e2074-1056-4b27-b4e0-ed9f0c46b067",
			err:       errors.New("{\"error_message\":\"invalid version number\"}"),
			code:      http.StatusBadRequest,
		},
		{
			name:      "happy path, account deleted",
			accountId: "cb1e2074-1056-4b27-b4e0-ed9f0c46b067",
			version:   "0",
			err:       nil,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			client := form3_client.Form3Client{
				HttpClient: &http.Client{},
				BaseURL:    getEnv("BASE_URL", "http://localhost:8080/"),
			}

			if test.name == "happy path, account deleted" {
				createDummyAccount()
			}
			err := client.DeleteAccount(test.accountId, test.version)
			if err.Error != nil {
				assert.Equal(t, test.err.Error(), err.Error.Error())
				assert.Equal(t, test.code, err.Code)
			} else {
				deleteDummyAccount(test.accountId)
			}
		})
	}
}

func createDummyAccount() {
	client := form3_client.Form3Client{
		HttpClient: &http.Client{},
		BaseURL:    getEnv("BASE_URL", "http://localhost:8080/"),
	}
	body := getBody()
	client.PostAccount(body)
}

func getBody() io.Reader {
	byteBody := []byte("{\n    \"data\": {\n        \"attributes\": {\n            \"account_classification\": \"Personal\",\n            \"account_matching_opt_out\": false,\n            \"alternative_names\": [\n                \"Sam Holder\"\n            ],\n            \"bank_id\": \"400300\",\n            \"bank_id_code\": \"GBDSC\",\n            \"base_currency\": \"GBP\",\n            \"bic\": \"NWBKGB22\",\n            \"country\": \"GB\",\n            \"joint_account\": false,\n            \"name\": [\n                \"Samantha Holder\"\n            ],\n            \"secondary_identification\": \"A1B2C3D4\"\n        },\n        \"id\": \"cb1e2074-1056-4b27-b4e0-ed9f0c46b067\",\n        \"organisation_id\": \"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c\",\n        \"type\": \"accounts\",\n        \"version\": 0\n    }\n}")
	return bytes.NewReader(byteBody)
}

func deleteDummyAccount(accountId string) {

	client := form3_client.Form3Client{
		HttpClient: &http.Client{},
		BaseURL:    getEnv("BASE_URL", "http://localhost:8080/"),
	}

	client.DeleteAccount(accountId, "0")
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
