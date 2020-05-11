package handlers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Authenticate_ValidCredentials(t *testing.T) {
	t.Run("It returns JSON with login and jwt token", func(t *testing.T) {
		data := url.Values{"login": {"TestRssApp"}, "password": {"password"}}
		request, err := http.NewRequest("POST", "http://host/authenticate", strings.NewReader(data.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		if err != nil {
			t.Error(err)
		}

		handler := NewAuthHandler()

		response, err := serveHTTPRequest(handler.Authenticate, request)
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 200, response.StatusCode)
		assert.True(t, strings.Contains(string(body), "{\"login\":\"TestRssApp\",\"token\":\""))
	})
}

func TestAuthHandler_Authenticate_InvalidCredentials(t *testing.T) {
	t.Run("It returns JSON with error message", func(t *testing.T) {
		data := url.Values{"login": {"TestRssApp"}, "password": {"wrong_password"}}
		request, err := http.NewRequest("POST", "http://host/authenticate", strings.NewReader(data.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		if err != nil {
			t.Error(err)
		}

		handler := NewAuthHandler()

		response, err := serveHTTPRequest(handler.Authenticate, request)
		if err != nil {
			t.Error(err)
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 401, response.StatusCode)
		assert.True(t, strings.Contains(string(body), "{\"error\":\"NotAuthenticated - your credentials are incorrect\""))
	})
}
