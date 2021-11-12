package form3_client

import (
	"encoding/json"
	"form3-interview/models"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	pathUrl = "v1/organisation/accounts"
)

type Form3ClientIface interface {
	GetAccount(accountId string) (account models.AccountWrapper, err error)
	PostAccount(body io.Reader) (account models.AccountWrapper, err error)
	DeleteAccount(accountId string, version string) (err error)
	Do(req *http.Request) (*http.Response, error)
}

type Form3Client struct {
	HttpClient *http.Client
}

func BuildBaseUrl() string {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		return "http://localhost:8080/"
	}
	return baseUrl
}

func (c Form3Client) GetAccount(accountId string) (account models.AccountWrapper, err error) {

	var (
		resp *http.Response
		req  *http.Request
	)
	url := BuildBaseUrl()
	fullUrl := url + pathUrl + "/" + accountId

	if req, err = http.NewRequest("GET", fullUrl, nil); err != nil {
		return account, errors.Wrap(err, "Unable to create new http client request")
	}

	if resp, err = c.Do(req); err != nil {
		return account, errors.Wrap(err, "Error during http request execution")
	}
	defer resp.Body.Close()

	if err = validation(resp); err != nil {
		return account, errors.Wrap(err, "Validation error")
	}

	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return account, errors.Wrap(err, "unable to decode the account response from form3 client")
	}
	return
}

func (c Form3Client) PostAccount(body io.Reader) (account models.AccountWrapper, err error) {
	var (
		resp *http.Response
		req  *http.Request
	)

	url := BuildBaseUrl()
	fullUrl := url + pathUrl

	if req, err = http.NewRequest("POST", fullUrl, body); err != nil {
		return account, errors.Wrap(err, "Unable to create new http client request")
	}

	if resp, err = c.Do(req); err != nil {
		return account, errors.Wrap(err, "Error during http request execution")
	}
	defer resp.Body.Close()

	if err = validation(resp); err != nil {
		return account, errors.Wrap(err, "Validation error")
	}

	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return account, errors.Wrap(err, "unable to decode the response from form3 client")
	}
	return
}

func (c Form3Client) DeleteAccount(accountId string, version string) (err error) {

	var req *http.Request
	url := BuildBaseUrl()

	fullUrl := url + pathUrl + "/" + accountId + "?version=" + version

	if req, err = http.NewRequest("DELETE", fullUrl, nil); err != nil {
		return errors.Wrap(err, "Unable to create new http client request")
	}
	if _, err = c.Do(req); err != nil {
		return errors.Wrap(err, "Error during http request execution "+accountId+",account not found")
	}
	return
}

func validation(resp *http.Response) (err error) {

	status := resp.StatusCode
	if status == http.StatusOK || status == http.StatusNoContent || status == http.StatusCreated {
		return nil
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		err := errors.New(string(respBody))
		log.Println(err)
		return err
	}
}

/*func (c Form3Client) doForm3HttpRequest(url string, body io.Reader, method string) (*http.Response, error) {
	var (
		req *http.Request
		err error
	)
	if req, err = http.NewRequest(method, url, body); err != nil {
		return nil, errors.Wrap(err, "Unable to create new http client request")
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	//resp, err := client.Do(req)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//defer resp.Body.Close()
	status := resp.StatusCode
	if status == http.StatusOK || status == http.StatusNoContent || status == http.StatusCreated {
		return resp, nil
	} else {
		respBody, _ := ioutil.ReadAll(resp.Body)
		err := errors.New(string(respBody))
		log.Println(err)
		return nil, err
	}
}*/

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
