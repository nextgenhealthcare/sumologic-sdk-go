package sumologic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultAWSCloudTrailSource = AWSCloudTrailSource{
	ID:          1234567890,
	Name:        "test",
	CollectorID: 1234567890,
}

func TestGetAWSCloudTrailSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(AWSCloudTrailSourceRequest{
			Source: defaultAWSCloudTrailSource,
		})
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedSource, _, err := c.GetAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
	if err != nil {
		t.Errorf("GetAWSCloudTrailSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != defaultAWSCloudTrailSource.ID {
		t.Errorf("GetAWSCloudTrailSource() expected ID `%d`, got `%d`", defaultAWSCloudTrailSource.ID, returnedSource.ID)
		return
	}
}

func TestGetAWSCloudTrailSourceDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
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

	_, _, err = c.GetAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
	if err != ErrSourceNotFound {
		t.Errorf("GetAWSCloudTrailSource() returned the wrong error: %s", err)
		return
	}
}

func TestCreateAWSCloudTrailSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources", defaultAWSCloudTrailSource.CollectorID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘/collectors’, got ‘%s’", r.URL.EscapedPath())
		}
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			t.Errorf("Expected response to be content-type ‘application/json’, got ‘%s’", ctype)
		}
		body, _ := ioutil.ReadAll(r.Body)
		sr := new(AWSCloudTrailSourceRequest)
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

	returnedSource, err := c.CreateAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, AWSCloudTrailSource{
		Name: "test",
	})
	if err != nil {
		t.Errorf("CreateAWSCloudTrailSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != 1234567890 {
		t.Errorf("CreateAWSCloudTrailSource() expected ID 1234567890, got `%d`", returnedSource.ID)
		return
	}
}

func TestCreateAWSCloudTrailSourceAlreadyExists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources", defaultAWSCloudTrailSource.CollectorID)
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

	_, err = c.CreateAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, AWSCloudTrailSource{
		Name: "test",
	})
	if err == nil {
		t.Errorf("CreateAWSCloudTrailSource() did not return an error: %s", err)
		return
	}
}

func TestUpdateAWSCloudTrailSourceOK(t *testing.T) {
	updatedSource := defaultAWSCloudTrailSource
	updatedSource.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		if r.Header.Get("If-Match") != "etag" {
			t.Errorf("Expected Etag of `etag`, got `%s`", r.Header.Get("If-Match"))
		}
		body, _ := json.Marshal(AWSCloudTrailSourceRequest{
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

	returnedSource, err := c.UpdateAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, updatedSource, "etag")
	if err != nil {
		t.Errorf("UpdateAWSCloudTrailSource() returned an error: %s", err)
		return
	}
	if returnedSource.ID != updatedSource.ID {
		t.Errorf("UpdateAWSCloudTrailSource() expected ID `%d`, got `%d`", defaultAWSCloudTrailSource.ID, returnedSource.ID)
		return
	}
	if returnedSource.Name == defaultAWSCloudTrailSource.Name {
		t.Errorf("UpdateAWSCloudTrailSource() did not update the name")
		return
	}
}

func TestUpdateAWSCloudTrailSourceAlreadyExists(t *testing.T) {
	updatedSource := defaultAWSCloudTrailSource
	updatedSource.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
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

	_, err = c.UpdateAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, updatedSource, "etag")
	if err == nil {
		t.Errorf("UpdateAWSCloudTrailSource() did not return an error: %s", err)
		return
	}
}

func TestDeleteAWSCloudTrailSourceOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
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

	err = c.DeleteAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
	if err != nil {
		t.Errorf("DeleteAWSCloudTrailSource() returned an error: %s", err)
		return
	}
}

func TestDeleteAWSCloudTrailSourceDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d/sources/%d", defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
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

	err = c.DeleteAWSCloudTrailSource(defaultAWSCloudTrailSource.CollectorID, defaultAWSCloudTrailSource.ID)
	if err != ErrSourceNotFound {
		t.Errorf("DeleteAWSCloudTrailSource() returned the wrong error: %s", err)
		return
	}
}
