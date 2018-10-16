package sumologic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultAWSLogSource = AWSLogSource{
	ID:          1234567890,
	Name:        "test",
	CollectorID: 1234567890,
}

func TestGetAWSLogSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(AWSLogSourceRequest{
			Source: defaultAWSLogSource,
		})
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedSource, _, err := c.GetAWSLogSource(defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
	if err != nil {
		t.Errorf("GetAWSLogSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != defaultAWSLogSource.ID {
		t.Errorf("GetAWSLogSource() expected ID `%d`, got `%d`", defaultAWSLogSource.ID, returnedSource.ID)
		return
	}
}

func TestGetAWSLogSourceDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
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

	_, _, err = c.GetAWSLogSource(defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
	if err != ErrSourceNotFound {
		t.Errorf("GetAWSLogSource() returned the wrong error: %s", err)
		return
	}
}

func TestCreateAWSLogSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources", defaultAWSLogSource.CollectorID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘/collectors’, got ‘%s’", r.URL.EscapedPath())
		}
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			t.Errorf("Expected response to be content-type ‘application/json’, got ‘%s’", ctype)
		}
		body, _ := ioutil.ReadAll(r.Body)
		sr := new(AWSLogSourceRequest)
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

	returnedSource, err := c.CreateAWSLogSource(defaultAWSLogSource.CollectorID, AWSLogSource{
		Name: "test",
	})
	if err != nil {
		t.Errorf("CreateAWSLogSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != 1234567890 {
		t.Errorf("CreateAWSLogSource() expected ID 1234567890, got `%d`", returnedSource.ID)
		return
	}
}

func TestCreateAWSLogSourceAlreadyExists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources", defaultAWSLogSource.CollectorID)
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

	_, err = c.CreateAWSLogSource(defaultAWSLogSource.CollectorID, AWSLogSource{
		Name: "test",
	})
	if err == nil {
		t.Errorf("CreateAWSLogSource() did not return an error: %s", err)
		return
	}
}

func TestUpdateAWSLogSourceOK(t *testing.T) {
	updatedSource := defaultAWSLogSource
	updatedSource.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		if r.Header.Get("If-Match") != "etag" {
			t.Errorf("Expected Etag of `etag`, got `%s`", r.Header.Get("If-Match"))
		}
		body, _ := json.Marshal(AWSLogSourceRequest{
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

	returnedSource, err := c.UpdateAWSLogSource(defaultAWSLogSource.CollectorID, updatedSource, "etag")
	if err != nil {
		t.Errorf("UpdateAWSLogSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != updatedSource.ID {
		t.Errorf("UpdateAWSLogSource() expected ID `%d`, got `%d`", defaultAWSLogSource.ID, returnedSource.ID)
		return
	}
	if returnedSource.Name == defaultAWSLogSource.Name {
		t.Errorf("UpdateAWSLogSource() did not update the name")
		return
	}
}

func TestUpdateAWSLogSourceAlreadyExists(t *testing.T) {
	updatedSource := defaultAWSLogSource
	updatedSource.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
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

	_, err = c.UpdateAWSLogSource(defaultAWSLogSource.CollectorID, updatedSource, "etag")
	if err == nil {
		t.Errorf("UpdateAWSLogSource() did not return an error: %s", err)
		return
	}
}

func TestDeleteAWSLogSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
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

	err = c.DeleteAWSLogSource(defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
	if err != nil {
		t.Errorf("DeleteAWSLogSource() returned an error: %s", err)
		return
	}
}

func TestDeleteAWSLogSourceDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
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

	err = c.DeleteAWSLogSource(defaultAWSLogSource.CollectorID, defaultAWSLogSource.ID)
	if err != ErrSourceNotFound {
		t.Errorf("DeleteAWSLogSource() returned the wrong error: %s", err)
		return
	}
}
