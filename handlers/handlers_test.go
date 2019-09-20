package handlers

import (
	"bytes"
	"encoding/json"
	"golang-fifa-world-cup-web-service/data"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"testing"
)

// reloads JSON into memory to ensure
// proper winner count during tests.
func setup() {
	p, _ := filepath.Abs("./../data/")
	fullpath := path.Join(p, "winners.json")
	data.LoadFromJSON(fullpath)
}

func TestRootHandlerReturnsNoContentStatus(t *testing.T) {
	handler := http.HandlerFunc(RootHandler)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Error("Did not return status 204 - No Content")
	}
}

func TestListWinnersSetsContentType(t *testing.T) {
	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	handler.ServeHTTP(rr, req)
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Error("Did not set Content-Type response header to application/json")
	}
}

func TestListWinnersReturnsAllWinners(t *testing.T) {
	setup()

	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	var winners data.Winners
	json.Unmarshal([]byte(body), &winners)
	if len(winners.Winners) != 21 {
		t.Error("Did not return all winners from /winners")
	}
}

func TestListWinnersReturnsAllWinnersByYear(t *testing.T) {
	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	q := req.URL.Query()
	q.Add("year", "2018")
	req.URL.RawQuery = q.Encode()

	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	var winners data.Winners
	json.Unmarshal([]byte(body), &winners)
	if len(winners.Winners) != 1 {
		t.Error("Did not return winners filtered by year")
	}
}

func TestListWinnersReturnsBadRequestWhenInvalidYear(t *testing.T) {
	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	q := req.URL.Query()
	q.Add("year", "banana")
	req.URL.RawQuery = q.Encode()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Error("Did not return status 400 - Bad Request")
	}
}

func TestAddNewWinnerHandlerReturnsUnauthorizedForInvalidAccessToken(t *testing.T) {
	setup()

	req, _ := http.NewRequest("POST", "/winners", nil)
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken+"bla")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Error("Did not return status 401 - Unauthorized for invalid Access Token")
	}
}

func TestAddNewWinnerHandlerAddsNewWinnerWithValidData(t *testing.T) {
	setup()

	var jsonStr = []byte(`{"country":"Croatia", "year": 2030}`)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)

	allWinners, _ := data.ListAllJSON()
	var winners data.Winners
	json.Unmarshal([]byte(allWinners), &winners)

	if len(winners.Winners) != 22 {
		t.Error("Did not properly add new winner to the list")
	}
}

func TestAddNewWinnerHandlerReturnsUnprocessableEntityForEmptyPayload(t *testing.T) {
	setup()

	// Invalid because empty
	var jsonStr = []byte(``)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Error("Did not properly validate winner payload")
	}
}

func TestAddNewWinnerHandlerDoesNotAddInvalidNewWinner(t *testing.T) {
	setup()

	// Invalid entry because year is in the past.
	var jsonStr = []byte(`{"country":"Croatia", "year": 1984}`)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddNewWinner)
	handler.ServeHTTP(rr, req)

	allWinners, _ := data.ListAllJSON()
	var winners data.Winners
	json.Unmarshal([]byte(allWinners), &winners)

	if rr.Code == http.StatusOK || len(winners.Winners) != 21 {
		t.Error("Added invalid winner to list")
	}
}

func TestCorrectHTTPGetMethodDispatch(t *testing.T) {
	setup()

	req, _ := http.NewRequest("GET", "/winners", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WinnersHandler)
	handler.ServeHTTP(rr, req)

	if body := rr.Body.String(); body == "" {
		t.Error("Did not properly dispatch HTTP GET")
	}
}

func TestCorrectHTTPPostMethodDispatch(t *testing.T) {
	setup()

	var jsonStr = []byte(`{"country":"Croatia", "year": 2030}`)
	req, _ := http.NewRequest("POST", "/winners", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-ACCESS-TOKEN", data.AccessToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WinnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Error("Did not properly dispatch HTTP POST", status)
	}
}

func TestCorrectHTTPUnsupportedMethodDispatch(t *testing.T) {
	setup()

	req, _ := http.NewRequest("PUT", "/winners", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WinnersHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Error("Did not properly catch unsupported HTTP methods", status)
	}
}
