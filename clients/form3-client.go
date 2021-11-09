package form3_client

import (
	"encoding/json"
	"form3-interview/models"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	pathUrl = "v1/organisation/accounts"
)

type Form3ClientIface interface {
	GetAccount(accountId string) (account models.AccountWrapper, err error)
	PostAccount(body io.Reader) (account models.AccountWrapper, err error)
}

type Form3Client struct {
	HttpClient *http.Client
	RetryTimeout time.Duration
}

func BuildBaseUrl() string {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == ""{
		return "http://localhost:8080/"
	}
	return baseUrl
}

func (c Form3Client) GetAccount(accountId string) (account models.AccountWrapper, err error) {
	var resp *http.Response
	url := BuildBaseUrl()
	fullUrl := url+pathUrl+"/"+accountId

	if resp, err = doForm3HttpRequest(fullUrl, nil, "GET"); err != nil {
		return account, errors.Wrap(err, "Unable to create get request for account")
	}
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return account, errors.Wrap(err, "unable to decode the account response from form3 client")
	}
	return
}

func (c Form3Client) PostAccount(body io.Reader) (account models.AccountWrapper, err error) {
	var resp *http.Response
	url := BuildBaseUrl()
	fullUrl := url+pathUrl

	if resp, err = doForm3HttpRequest(fullUrl, body , "POST"); err != nil {
		return account, errors.Wrap(err, "Unable to create post request for account")
	}
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return account, errors.Wrap(err, "unable to decode the response from form3 client")
	}
	return
}

func doForm3HttpRequest(url string, body io.Reader, method string) (*http.Response, error) {
	var (
		req         *http.Request
		err         error
	)
	if req, err = http.NewRequest(method, url, body); err != nil {
		return nil, errors.Wrap(err, "Unable to create new http client request")
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
