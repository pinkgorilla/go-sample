package response_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pinkgorilla/go-sample/starter-api/errors"
	"github.com/pinkgorilla/go-sample/starter-api/internal/http-server/response"
)

func getServer(err error) *httptest.Server {
	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response.WithError(w, nil, err)
		}))
	return s
}

func Test_Response_WithError(t *testing.T) {
	cases := []struct {
		err      error // error type to be tested
		expected int   // expected http status code}
	}{
		{errors.NewBaseError("code", "error"), http.StatusBadRequest},
		{errors.NewCommonError("error"), http.StatusBadRequest},
		{errors.NewAuthError("error"), http.StatusUnauthorized},
		{errors.NewPermissionError("error"), http.StatusForbidden},
		{errors.NewNotFoundError("error"), http.StatusNotFound},
		{errors.NewValidationError("error"), http.StatusUnprocessableEntity},
		{errors.NewServiceError("error"), http.StatusInternalServerError},
		{fmt.Errorf("error"), http.StatusInternalServerError},
	}

	for _, c := range cases {
		s := getServer(c.err)
		res, err := s.Client().Get(s.URL)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != c.expected {
			t.Fatalf("expected status %v, get %v", c.expected, res.StatusCode)
		}

		var d map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&d)
		if err != nil {
			t.Fatal(err)
		}
		_, ok := d["code"]
		if !ok {
			t.Fatal("expected 'code' key in error response")
		}
		_, ok = d["message"]
		if !ok {
			t.Fatal("expected 'message' key in error response")
		}

		s.Close()
	}

}
