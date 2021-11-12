package handlers_test

import (
	"form3-interview/handlers"
	mock_form3_client "form3-interview/mocks"
	"form3-interview/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pariz/gountries"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_form3GetHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		givenPayload interface{}
		err          error
		pathParam    map[string]string
		mockShop     func(mock *mock_form3_client.MockForm3ClientIface)
		status       int
	}{
		{
			name:      "accountId, not provided",
			pathParam: nil,
			mockShop:  func(mock *mock_form3_client.MockForm3ClientIface) {},
			status:    http.StatusBadRequest,
		},
		{
			name:      "happy path, accounts created",
			pathParam: map[string]string{"accountId": ""},
			mockShop: func(mock *mock_form3_client.MockForm3ClientIface) {
				mock.EXPECT().GetAccount(gomock.Any()).Return(mockedAccount(), nil)
			},
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", "/form3Client/accounts/", nil)
			if err != nil {
				t.Errorf("Error creating a new request: %v", err)
			}
			req = mux.SetURLVars(req, test.pathParam)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_form3_client.NewMockForm3ClientIface(ctrl)
			rr := httptest.NewRecorder()
			test.mockShop(mockClient)
			handler := http.HandlerFunc(handlers.GetAccount(mockClient))
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
		})
	}
}

func Test_form3DeleteHandler(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		pathParam map[string]string
		version   string
		mockShop  func(mock *mock_form3_client.MockForm3ClientIface)
		status    int
	}{
		{
			name:      "accountId, not provided",
			pathParam: nil,
			version:   "nil",
			mockShop:  func(mock *mock_form3_client.MockForm3ClientIface) {},
			status:    http.StatusBadRequest,
		},
		{
			name:      "version, not provided",
			pathParam: map[string]string{"accountId": ""},
			version:   "",
			mockShop:  func(mock *mock_form3_client.MockForm3ClientIface) {},
			status:    http.StatusBadRequest,
		},
		{
			name:      "happy path, deleted",
			pathParam: map[string]string{"accountId": "1234"},
			version:   "1",
			mockShop: func(mock *mock_form3_client.MockForm3ClientIface) {
				mock.EXPECT().DeleteAccount(gomock.Any(), gomock.Any()).Return(nil)
			},
			status: http.StatusNoContent,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("DELETE", "/form3Client/accounts/", nil)
			if err != nil {
				t.Errorf("Error creating a new request: %v", err)
			}
			req = mux.SetURLVars(req, test.pathParam)
			q := req.URL.Query()
			q.Add("version", test.version)
			req.URL.RawQuery = q.Encode()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_form3_client.NewMockForm3ClientIface(ctrl)
			rr := httptest.NewRecorder()
			test.mockShop(mockClient)
			handler := http.HandlerFunc(handlers.DeleteAccount(mockClient))
			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
		})
	}
}

func mockedAccount() models.AccountWrapper {

	query := gountries.New()
	sweden, _ := query.FindCountryByName("canada")

	return models.AccountWrapper{
		Account: models.AccountData{
			ID: "60c6add9-2b7b-4427-972a-8b272735562f",
			Attributes: &models.AccountAttributes{
				BankID:  "123456",
				Country: &sweden.Alpha2,
			},
		},
	}
}
