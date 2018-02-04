package sumologic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var defaultCollector = Collector{
	ID:            1234567890,
	Name:          "test",
	CollectorType: "Hosted",
}

func TestAuthenticationFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d", defaultCollector.ID)
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

	_, _, err = c.GetHostedCollector(defaultCollector.ID)
	if err != ErrClientAuthenticationError {
		t.Errorf("GetHostedCollector() returned the wrong error: %s", err)
		return
	}
}

func TestGetHostedCollectorOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d", defaultCollector.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		body, _ := json.Marshal(CollectorRequest{
			Collector: defaultCollector,
		})
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCollector, _, err := c.GetHostedCollector(defaultCollector.ID)
	if err != nil {
		t.Errorf("GetHostedCollector() returned an error: %s", err)
		return
	}
	if returnedCollector.ID != defaultCollector.ID {
		t.Errorf("GetHostedCollector() expected ID `%d`, got `%d`", defaultCollector.ID, returnedCollector.ID)
		return
	}
}

func TestGetHostedCollectorDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "GET" {
			t.Errorf("Expected ‘GET’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d", defaultCollector.ID)
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

	_, _, err = c.GetHostedCollector(defaultCollector.ID)
	if err != ErrCollectorNotFound {
		t.Errorf("GetHostedCollector() returned the wrong error: %s", err)
		return
	}
}

func TestCreateHostedCollectorOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		if r.URL.EscapedPath() != "/collectors" {
			t.Errorf("Expected request to ‘/collectors’, got ‘%s’", r.URL.EscapedPath())
		}
		if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
			t.Errorf("Expected response to be content-type ‘application/json’, got ‘%s’", ctype)
		}
		body, _ := ioutil.ReadAll(r.Body)
		cr := new(CollectorRequest)
		err := json.Unmarshal(body, &cr)
		if err != nil {
			t.Errorf("Unable to unmarshal CollectorRequest, got `%s`", body)
		}
		if cr.Collector.Name != "test" {
			t.Errorf("Expected request to include collector name ‘test’, got ‘%s’", cr.Collector.Name)
		}
		cr.Collector.ID = 1234567890
		js, err := json.Marshal(cr)
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

	returnedCollector, err := c.CreateHostedCollector(Collector{
		Name:          "test",
		CollectorType: "Hosted",
	})
	if err != nil {
		t.Errorf("CreateHostedCollector() returned an error: %s", err)
		return
	}
	if returnedCollector.ID != 1234567890 {
		t.Errorf("CreateHostedCollector() expected ID 1234567890, got `%d`", returnedCollector.ID)
		return
	}
}

func TestCreateHostedCollectorAlreadyExists(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}
		if r.URL.EscapedPath() != "/collectors" {
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

	_, err = c.CreateHostedCollector(Collector{
		Name:          "test",
		CollectorType: "Hosted",
	})
	if err == nil {
		t.Errorf("CreateHostedCollector() did not return an error: %s", err)
		return
	}
}

func TestUpdateHostedCollectorOK(t *testing.T) {
	updatedCollector := defaultCollector
	updatedCollector.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d", defaultCollector.ID)
		if r.URL.EscapedPath() != expectedURL {
			t.Errorf("Expected request to ‘%s’, got ‘%s’", expectedURL, r.URL.EscapedPath())
		}
		if r.Header.Get("If-Match") != "etag" {
			t.Errorf("Expected Etag of `etag`, got `%s`", r.Header.Get("If-Match"))
		}
		body, _ := json.Marshal(CollectorRequest{
			Collector: updatedCollector,
		})
		w.Write(body)
	}))
	defer ts.Close()

	c, err := NewClient("accessToken", ts.URL)
	if err != nil {
		t.Errorf("NewClient() returned an error: %s", err)
		return
	}

	returnedCollector, err := c.UpdateHostedCollector(updatedCollector, "etag")
	if err != nil {
		t.Errorf("UpdateHostedCollector() returned an error: %s", err)
		return
	}
	if returnedCollector.ID != updatedCollector.ID {
		t.Errorf("UpdateHostedCollector() expected ID `%d`, got `%d`", defaultCollector.ID, returnedCollector.ID)
		return
	}
	if returnedCollector.Name == defaultCollector.Name {
		t.Errorf("UpdateHostedCollector() did not update the name")
		return
	}
}

func TestUpdateHostedCollectorAlreadyExists(t *testing.T) {
	updatedCollector := defaultCollector
	updatedCollector.Name = "Updated"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if r.Method != "PUT" {
			t.Errorf("Expected ‘PUT’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d", defaultCollector.ID)
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

	_, err = c.UpdateHostedCollector(updatedCollector, "etag")
	if err == nil {
		t.Errorf("UpdateHostedCollector() did not return an error: %s", err)
		return
	}
}

func TestDeleteHostedCollectorOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d", defaultCollector.ID)
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

	err = c.DeleteHostedCollector(defaultCollector.ID)
	if err != nil {
		t.Errorf("DeleteHostedCollector() returned an error: %s", err)
		return
	}
}

func TestDeleteHostedCollectorDoesntExist(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if r.Method != "DELETE" {
			t.Errorf("Expected ‘DELETE’ request, got ‘%s’", r.Method)
		}
		expectedURL := fmt.Sprintf("/collectors/%d", defaultCollector.ID)
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

	err = c.DeleteHostedCollector(defaultCollector.ID)
	if err != ErrCollectorNotFound {
		t.Errorf("DeleteHostedCollector() returned the wrong error: %s", err)
		return
	}
}
