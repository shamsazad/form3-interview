package form3_client_test

import (
	"bytes"
	"encoding/json"
	form3_client "form3-interview/clients"
	"form3-interview/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_form3ClientGet(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		givenPayload interface{}
		err          error
		pathParam    string
		testServer   *httptest.Server
		status       int
		separator    string
	}{
		{
			name:      "broken url for GET request",
			pathParam: "",

			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(500)
			})),
			status:    http.StatusInternalServerError,
			err:       errors.New("Malfunctioned http client request"),
			separator: "",
		},
		{
			name:      "unable to reach server",
			pathParam: "",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(500)
			})),
			status:    http.StatusInternalServerError,
			err:       errors.New("Unable to reach form3 server"),
			separator: "/",
		},
		{
			name:      "bad data coming from server",
			pathParam: "cb1e2074-1056-4b27-b4e0-ed9f0c46b067",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
				res.Write([]byte("{data {bad:data}"))
			})),
			status:    http.StatusInternalServerError,
			err:       errors.New("Unable to decode the account response from form3 client"),
			separator: "/",
		},
		{
			name:      "Account id doesn't exist",
			pathParam: "cb1e2074-1056-4b27-b4e0-ed9f0c46b067",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(400)
				res.Write([]byte("{error_message:record cb1e2074-1056-4b27-b4e0-ed9f0c46b067 does not exist}"))
			})),
			status:    http.StatusInternalServerError,
			err:       errors.New("Validation error"),
			separator: "/",
		},
		{
			name:      "happy path, account retrieved",
			pathParam: "",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
				res.Write(createDummyAccount())
			})),
			status:    http.StatusOK,
			err:       nil,
			separator: "/",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testServer := test.testServer
			defer testServer.Close()
			if test.name == "unable to reach server" {
				testServer.Close()
			}

			client := form3_client.Form3Client{
				HttpClient: testServer.Client(),
				BaseURL:    testServer.URL + test.separator,
			}

			account, err := client.GetAccount(test.pathParam)
			if test.err == nil {
				var expectedAccount models.AccountWrapper
				err := json.Unmarshal(createDummyAccount(), &expectedAccount)
				if err != nil {
					t.Fatalf("unable to marshal account")
				}
				assert.Equal(t, account, expectedAccount)
			} else {
				//v := strings.Split(errors.Unwrap(err).Error(), ":")
				assert.Equal(t, test.err.Error(), err.Message)
			}
		})
	}
}

func Test_form3ClientPost(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		givenPayload []byte
		err          error
		testServer   *httptest.Server
		separator    string
	}{
		{
			name:         "broken url for Post request",
			givenPayload: []byte(""),
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			})),
			err:       errors.New("Malfunctioned http client request"),
			separator: "",
		},
		{
			name:         "unable to reach server",
			givenPayload: []byte(""),
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			})),
			err:       errors.New("Unable to reach form3 server"),
			separator: "/",
		},
		{
			name:         "bad data coming from server",
			givenPayload: []byte(""),
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
				res.Write([]byte("{data {bad:data}"))
			})),
			err:       errors.New("Unable to decode the account response from form3 client"),
			separator: "/",
		},
		{
			name:         "Missing required field, organization_id",
			givenPayload: []byte(""),
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(400)
				res.Write([]byte("{error_message:validation failure list:\\nvalidation failure list:\\norganisation_id in body is required}"))
			})),
			err:       errors.New("Validation error"),
			separator: "/",
		},
		{
			name:         "happy path, account created",
			givenPayload: createDummyAccount(),
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(200)
				res.Write(createDummyAccount())
			})),
			err:       nil,
			separator: "/",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testServer := test.testServer
			defer testServer.Close()
			if test.name == "unable to reach server" {
				testServer.Close()
			}

			client := form3_client.Form3Client{
				HttpClient: testServer.Client(),
				BaseURL:    testServer.URL + test.separator,
			}

			body := bytes.NewReader(test.givenPayload)
			account, err := client.PostAccount(body)
			if test.err == nil {
				var expectedAccount models.AccountWrapper
				err := json.Unmarshal(createDummyAccount(), &expectedAccount)
				if err != nil {
					t.Fatalf("unable to marshal account")
				}
				assert.Equal(t, account, expectedAccount)
			} else {
				//v := strings.Split(errors.Unwrap(err).Error(), ":")
				assert.Equal(t, test.err.Error(), err.Message)
			}
		})
	}
}

func Test_form3ClientDelete(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		accountId  string
		version    string
		err        error
		testServer *httptest.Server
		separator  string
	}{
		{
			name: "broken url for delete request",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			})),
			err:       errors.New("Malfunctioned http client request"),
			separator: "",
		},
		{
			name: "unable to reach server",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			})),
			err:       errors.New("Unable to reach form3 server"),
			separator: "/",
		},
		{
			name: "Missing required field, account_id",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(400)
				res.Write([]byte("{error_message:validation failure list:\\nvalidation failure list:\\norganisation_id in body is required}"))
			})),
			err:       errors.New("Validation error"),
			separator: "/",
		},
		{
			name: "Missing required field, version",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(400)
				res.Write([]byte("{error_message:validation failure list:\\nvalidation failure list:\\norganisation_id in body is required}"))
			})),
			err:       errors.New("Validation error"),
			separator: "/",
		},
		{
			name: "happy path, account deleted",
			testServer: httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(201)
				res.Write(createDummyAccount())
			})),
			err:       nil,
			separator: "/",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			testServer := test.testServer
			defer testServer.Close()
			if test.name == "unable to reach server" {
				testServer.Close()
			}

			client := form3_client.Form3Client{
				HttpClient: testServer.Client(),
				BaseURL:    testServer.URL + test.separator,
			}

			err := client.DeleteAccount(test.accountId, test.version)
			if err.Error != nil {
				assert.Equal(t, test.err.Error(), err.Message)
			}
		})
	}
}

func createDummyAccount() []byte {
	return []byte("{\n    \"data\": {\n        \"attributes\": {\n            \"account_classification\": \"Personal\",\n            \"account_matching_opt_out\": false,\n            \"alternative_names\": [\n                \"Sam Holder\"\n            ],\n            \"bank_id\": \"400300\",\n            \"bank_id_code\": \"GBDSC\",\n            \"base_currency\": \"GBP\",\n            \"bic\": \"NWBKGB22\",\n            \"country\": \"GB\",\n            \"joint_account\": false,\n            \"name\": [\n                \"Samantha Holder\"\n            ],\n            \"secondary_identification\": \"A1B2C3D4\"\n        },\n        \"id\": \"cb1e2074-1056-4b27-b4e0-ed9f0c46b066\",\n        \"organisation_id\": \"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c\",\n        \"type\": \"accounts\",\n        \"version\": 0\n    }\n}")
}
