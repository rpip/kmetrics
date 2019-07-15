package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHealthCheckHandler(t *testing.T) {
	t.Parallel()

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctx := CreateContextForTestSetup()
	handler := makeHandler(ctx, HealthCheckHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"appName":"KMetrics","version":"0.0.0"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListServicesHandler(t *testing.T) {
	t.Parallel()

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/services", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	ctx := CreateContextForTestSetup()
	handler := makeHandler(ctx, ListServicesHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	data, _ := ioutil.ReadFile("./testdata/services.json")
	expected := stripSpaces(string(data))
	//strings.Replace(string(data), " ", "", -1)

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetServiceGroupHandler(t *testing.T) {
	t.Parallel()

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/services", nil)
	if err != nil {
		t.Fatal(err)
	}
	// set route named params
	vars := map[string]string{"group": "beta"}
	req = mux.SetURLVars(req, vars)

	rr := httptest.NewRecorder()
	ctx := CreateContextForTestSetup()
	handler := makeHandler(ctx, SearchServicesHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	data, _ := ioutil.ReadFile("./testdata/beta.json")
	expected := stripSpaces(string(data))
	//strings.Replace(string(data), " ", "", -1)

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
