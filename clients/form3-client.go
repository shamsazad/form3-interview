package form3_client

import (
	"encoding/json"
	"form3-interview/models"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	pathUrl = "v1/organisation/accounts"
)

type Form3ClientIface interface {
	GetAccount(accountId string) (account models.AccountWrapper, err models.AppError)
	PostAccount(body io.Reader) (account models.AccountWrapper, err models.AppError)
	DeleteAccount(accountId string, version string) (err models.AppError)
	Do(req *http.Request) (*http.Response, error)
}

type Form3Client struct {
	HttpClient *http.Client
	BaseURL    string
}

/*func BuildBaseUrl() string {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return "http://localhost:8080/"
	}
	return baseUrl
}*/

func (c Form3Client) GetAccount(accountId string) (account models.AccountWrapper, appError models.AppError) {

	var (
		resp *http.Response
		req  *http.Request
		err  error
	)
	url := c.BaseURL
	fullUrl := url + pathUrl + "/" + accountId

	if req, err = http.NewRequest("GET", fullUrl, nil); err != nil {
		return account, models.NewAppError(err, "Malfunctioned http client request", 500)
	}

	if resp, err = c.Do(req); err != nil {
		return account, models.NewAppError(err, "Unable to reach form3 server", 500)
	}
	defer resp.Body.Close()

	if appError = validation(resp); appError.Error != nil {
		return account, models.NewAppError(appError.Error, "Validation error", appError.Code)
	}

	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return account, models.NewAppError(err, "Unable to decode the account response from form3 client", 500)
	}
	return
}

func (c Form3Client) PostAccount(body io.Reader) (account models.AccountWrapper, appError models.AppError) {
	var (
		resp *http.Response
		req  *http.Request
		err  error
	)

	url := c.BaseURL
	fullUrl := url + pathUrl

	if req, err = http.NewRequest("POST", fullUrl, body); err != nil {
		return account, models.NewAppError(err, "Malfunctioned http client request", 500)
	}

	if resp, err = c.Do(req); err != nil {
		return account, models.NewAppError(err, "Unable to reach form3 server", 500)
	}
	defer resp.Body.Close()

	if appError = validation(resp); appError.Error != nil {
		return account, models.NewAppError(appError.Error, "Validation error", appError.Code)
	}

	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return account, models.NewAppError(appError.Error, "Unable to decode the account response from form3 client", 500)
	}
	return
}

func (c Form3Client) DeleteAccount(accountId string, version string) (appError models.AppError) {

	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	url := c.BaseURL

	fullUrl := url + pathUrl + "/" + accountId + "?version=" + version

	if req, err = http.NewRequest("DELETE", fullUrl, nil); err != nil {
		return models.NewAppError(err, "Malfunctioned http client request", 500)
	}
	if resp, err = c.Do(req); err != nil {
		return models.NewAppError(err, "Unable to reach form3 server", 500)
	}
	if appError = validation(resp); appError.Error != nil {
		return models.NewAppError(appError.Error, "Validation error", appError.Code)
	}
	return
}

func validation(resp *http.Response) (appError models.AppError) {

	status := resp.StatusCode
	if status == http.StatusOK || status == http.StatusNoContent || status == http.StatusCreated {
		return appError
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		err := errors.New(string(respBody))
		log.Println(err)
		return models.NewAppError(err, string(respBody), status)
	}
}

func (c *Form3Client) Do(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	req.Header.Set("Content-Type", "application/json")
	if resp, err = c.HttpClient.Do(req); err != nil {
		return nil, err
	}
	return resp, nil
}
