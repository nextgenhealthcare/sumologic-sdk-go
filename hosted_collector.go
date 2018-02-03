package sumologic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// CollectorRequest is a necessary wrapper for collector API calls.
type CollectorRequest struct {
	Collector Collector `json:"collector"`
}

// Collector can either be an installed or hosted collector.
// Installed collectors are installed as agents on servers.
// Hosted collectors receive data via HTTP or more specicialized (e.g. reading from AWS S3).
type Collector struct {
	ID               int              `json:"id,omitempty"`
	Name             string           `json:"name"`
	Description      string           `json:"description,omitempty"`
	Category         string           `json:"category,omitempty"`
	TimeZone         string           `json:"timezone,omitempty"`
	Links            []CollectorLinks `json:"links,omitempty"`
	CollectorType    string           `json:"collectorType,omitempty"`
	CollectorVersion string           `json:"collectorVersion,omitempty"`
	LastSeenAlive    int64            `json:"lastSeenAlive,omitempty"`
	Alive            bool             `json:"alive,omitempty"`
}

// CollectorLinks contains references to related resources such as sources.
type CollectorLinks struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// ErrCollectorNotFound is returned when a collector doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if exists).
var ErrCollectorNotFound = errors.New("Collector not found")

// GetHostedCollector gets the collector with the specified ID.
func (s *Client) GetHostedCollector(id int) (*Collector, string, error) {

	relativeURL, _ := url.Parse(fmt.Sprintf("collectors/%d", id))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("Authorization", "Basic "+s.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	ResponseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		var cr = new(CollectorRequest)
		err = json.Unmarshal(ResponseBody, &cr)
		if err != nil {
			return nil, "", err
		}

		return &cr.Collector, resp.Header.Get("ETag"), nil
	case http.StatusUnauthorized:
		return nil, "", ErrClientAuthenticationError
	case http.StatusNotFound:
		return nil, "", ErrCollectorNotFound
	default:
		return nil, "", fmt.Errorf("Unknown Response with Sumo Logic: `%s`", resp.StatusCode)
	}
}

// CreateHostedCollector creates a new Hosted Collector.
func (s *Client) CreateHostedCollector(collector Collector) (*Collector, error) {

	collectorRequest := CollectorRequest{
		Collector: collector,
	}

	body, _ := json.Marshal(collectorRequest)

	relativeURL, _ := url.Parse("collectors")
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+s.AuthToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusCreated:
		var cr = new(CollectorRequest)
		err = json.Unmarshal(responseBody, &cr)
		if err != nil {
			return nil, err
		}

		return &cr.Collector, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad Request. Please check if a collector with this name `%` already exists", collector.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with Sumo Logic: `%s`", resp.StatusCode)
	}
}

// UpdateHostedCollector updates an existing hosted collector.
func (s *Client) UpdateHostedCollector(collector Collector, etag string) (*Collector, error) {
	collectorRequest := CollectorRequest{
		Collector: collector,
	}

	body, _ := json.Marshal(collectorRequest)

	relativeURL, _ := url.Parse(fmt.Sprintf("collectors/%d", collector.ID))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("PUT", url.String(), bytes.NewBuffer((body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+s.AuthToken)
	req.Header.Add("If-Match", etag)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ResponseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		var cr = new(CollectorRequest)
		err = json.Unmarshal(ResponseBody, &cr)
		if err != nil {
			return nil, err
		}

		return &cr.Collector, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad Request. Please check if a collector with this name `%` already exists", collector.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with Sumo Logic: `%s`", resp.StatusCode)
	}
}

// DeleteHostedCollector deletes the collector with the specified ID.
func (s *Client) DeleteHostedCollector(id int) error {
	c, _ := url.Parse(fmt.Sprintf("collectors/%d", id))
	req, err := http.NewRequest("DELETE", s.EndpointURL.ResolveReference(c).String(), nil)
	req.Header.Add("Authorization", "Basic "+s.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return ErrCollectorNotFound
	case http.StatusUnauthorized:
		return ErrClientAuthenticationError
	default:
		return fmt.Errorf("Unknown Response with Sumo Logic: `%s`", resp.StatusCode)
	}
}
