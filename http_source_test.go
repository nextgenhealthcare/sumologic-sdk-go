package sumologic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultHTTPSource = HTTPSource{
	ID:          1234567890,
	Name:        "test",
	CollectorID: 1234567890,
}

func TestGetHTTPSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(HTTPSourceRequest{
			Source: defaultHTTPSource,
		})
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedSource, _, err := c.GetHTTPSource(defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
	if err != nil {
		t.Errorf("GetHTTPSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != defaultHTTPSource.ID {
		t.Errorf("GetHTTPSource() expected ID `%d`, got `%d`", defaultHTTPSource.ID, returnedSource.ID)
		return
	}
}

func TestGetHTTPSourceDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	_, _, err = c.GetHTTPSource(defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
	if err != ErrSourceNotFound {
		t.Errorf("GetHTTPSource() returned the wrong error: %s", err)
		return
	}
}

func TestCreateHTTPSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources", defaultHTTPSource.CollectorID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘/collectors’, got ‘%s’", r.URL.EscapedPath())
		}
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			t.Errorf("Expected response to be content-type ‘application/json’, got ‘%s’", ctype)
		}
		body, _ := ioutil.ReadAll(r.Body)
		sr := new(HTTPSourceRequest)
		err := json.Unmarshal(body, &sr)
		if err != nil {
			t.Errorf("Unable to unmarshal CollectorRequest, got `%s`", body)
		}
		if sr.Source.Name != "test" {
			t.Errorf("Expected request to include source name ‘test’, got ‘%s’", sr.Source.Name)
		}
		sr.Source.ID = 1234567890
		js, err := json.Marshal(sr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(js)
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedSource, err := c.CreateHTTPSource(defaultHTTPSource.CollectorID, HTTPSource{
		Name: "test",
	})
	if err != nil {
		t.Errorf("CreateHTTPSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != 1234567890 {
		t.Errorf("CreateHTTPSource() expected ID 1234567890, got `%d`", returnedSource.ID)
		return
	}
}

func TestCreateHTTPSourceAlreadyExists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources", defaultHTTPSource.CollectorID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘/collectors’, got ‘%s’", r.URL.EscapedPath())
		}
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			t.Errorf("Expected response to be content-type ‘application/json’, got ‘%s’", ctype)
		}
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	_, err = c.CreateHTTPSource(defaultHTTPSource.CollectorID, HTTPSource{
		Name: "test",
	})
	if err == nil {
		t.Errorf("CreateHTTPSource() did not return an error: %s", err)
		return
	}
}

func TestUpdateHTTPSourceOK(t *testing.T) {
	updatedSource := defaultHTTPSource
	updatedSource.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		if r.Header.Get("If-Match") != "etag" {
			t.Errorf("Expected Etag of `etag`, got `%s`", r.Header.Get("If-Match"))
		}
		body, _ := json.Marshal(HTTPSourceRequest{
			Source: updatedSource,
		})
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedSource, err := c.UpdateHTTPSource(defaultHTTPSource.CollectorID, updatedSource, "etag")
	if err != nil {
		t.Errorf("UpdateHTTPSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != updatedSource.ID {
		t.Errorf("UpdateHTTPSource() expected ID `%d`, got `%d`", defaultHTTPSource.ID, returnedSource.ID)
		return
	}
	if returnedSource.Name == defaultHTTPSource.Name {
		t.Errorf("UpdateHTTPSource() did not update the name")
		return
	}
}

func TestUpdateHTTPSourceAlreadyExists(t *testing.T) {
	updatedSource := defaultHTTPSource
	updatedSource.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		if r.Header.Get("If-Match") != "etag" {
			t.Errorf("Expected Etag of `etag`, got `%s`", r.Header.Get("If-Match"))
		}
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	_, err = c.UpdateHTTPSource(defaultHTTPSource.CollectorID, updatedSource, "etag")
	if err == nil {
		t.Errorf("UpdateHTTPSource() did not return an error: %s", err)
		return
	}
}

func TestDeleteHTTPSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	err = c.DeleteHTTPSource(defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
	if err != nil {
		t.Errorf("DeleteHTTPSource() returned an error: %s", err)
		return
	}
}

func TestDeleteHTTPSourceDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	err = c.DeleteHTTPSource(defaultHTTPSource.CollectorID, defaultHTTPSource.ID)
	if err != ErrSourceNotFound {
		t.Errorf("DeleteHTTPSource() returned the wrong error: %s", err)
		return
	}
}
