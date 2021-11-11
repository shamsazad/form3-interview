package handlers_test

import (
	"bytes"
	"encoding/json"
	form3_client "form3-interview/clients"
	"form3-interview/handlers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_form3Handler(t *testing.T) {
	t.Parallel()

	type testBed struct {
		request   *http.Request
		urlParams url.Values
		expectFn  func(response *httptest.ResponseRecorder, err error)
	}
	type testScenario struct {
		testName string
		testBed  func(t *testing.T) testBed
	}

	type buffaloError struct {
		Err string `json:"error"`
	}

	testScenarios := []testScenario{
		{
			testName: "Should return http.StatusBadRequest when `accountId` param is missing",
			testBed: func(t *testing.T) testBed {
				return testBed{
					request: httptest.NewRequest("", "/", strings.NewReader(`{`)),
					expectFn: func(response *httptest.ResponseRecorder, err error) {
						assert.NoError(t, err)
						var responseError buffaloError
						err = json.NewDecoder(response.Body).Decode(&responseError)
						assert.NoError(t, err)
						assert.Equal(t, 400, response.Result().StatusCode)
						assert.Contains(t, responseError.Err, "Missing `account` param")
					},
				}
			},
		},
	}

	for _, ts := range testScenarios {
		ts := ts
		t.Run(ts.testName, func(t *testing.T) {
			t.Parallel()
			Client := form3_client.Form3Client{}
			rr := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/posts", nil)
			handler := http.HandlerFunc(handlers.GetAccount(Client))
			handler.ServeHTTP(rr, req)
			//testBed := ts.testBed(t)

			//resp := httptest.NewRecorder()
			//err := handlers.GetAccount( nil)
			//testBed.expectFn(resp, err)
		})
	}
}

type JsonPlaceholderMock struct{}

type Post struct {
	Id     int
	UserId int
	Title  string
	Body   string
}

func (*JsonPlaceholderMock) GetPosts() (*http.Response, error) {
	mockedPosts := []Post{{}}

	respBody, err := json.Marshal(mockedPosts)

	if err != nil {
		log.Panicf("Error reading mocked response data: %v", err)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(respBody)),
	}, nil
}

func TestGetHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/form3Client/accounts/", nil)

	if err != nil {
		t.Errorf("Error creating a new request: %v", err)
	}

	m := map[string]string{
		"accountId": "60c6add9-2b7b-4427-972a-8b272735562f",
	}
	req = mux.SetURLVars(req, m)
	Client := form3_client.Form3Client{}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetAccount(Client))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var posts []Post

	if err := json.NewDecoder(rr.Body).Decode(&posts); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	resultTotal := len(posts)
	expectedTotal := 1

	if resultTotal != expectedTotal {
		t.Errorf("Expected: %d. Got: %d.", expectedTotal, resultTotal)
	}
}
