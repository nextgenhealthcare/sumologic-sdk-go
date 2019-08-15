package sumologic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// HTTPSource is a necessary wrapper for source API calls.
type HTTPSourceRequest struct {
	Source HTTPSource `json:"source"`
}

// HTTPSource can various types of sources including Cloudtrail and S3.
type HTTPSource struct {
	ID                         int      `json:"id,omitempty"`
	Name                       string   `json:"name"`
	CollectorID                int      `json:"CollectorId,omitempty"`
	Description                string   `json:"description,omitempty"`
	Category                   string   `json:"category,omitempty"`
	TimeZone                   string   `json:"timezone,omitempty"`
	SourceType                 string   `json:"sourceType,omitempty"`
	MessagePerRequest          bool     `json:"messagePerRequest"`
	MultilineProcessingEnabled bool     `json:"multilineProcessingEnabled"`
	UseAutolineMatching        bool     `json:"useAutolineMatching,"`
	ManualPrefixRegexp         string   `json:"manualPrefixRegexp,omitempty"`
	Url                        string   `json:"url,omitempty"`
	Filters                    []Filter `json:"filters,omitempty"`
}

// GetHTTPSource gets the source with the specified ID.
func (s *Client) GetHTTPSource(collectorID int, id int) (*HTTPSource, string, error) {

	relativeURL, _ := url.Parse(fmt.Sprintf("collectors/%d/sources/%d", collectorID, id))
	url := s.EndpointURL.ResolveReference(relativeURL)

	req, err := http.NewRequest("GET", url.String(), nil)
	req.Header.Add("Authorization", "Basic "+s.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		var r = new(HTTPSourceRequest)
		err = json.Unmarshal(responseBody, &r)
		if err != nil {
			return nil, "", err
		}

		return &r.Source, resp.Header.Get("ETag"), nil
	case http.StatusUnauthorized:
		return nil, "", ErrClientAuthenticationError
	case http.StatusNotFound:
		return nil, "", ErrSourceNotFound
	default:
		return nil, "", fmt.Errorf("Unknown Response with Sumo Logic: `%d`", resp.StatusCode)
	}
}

// CreateHTTPSource creates a new HTTPSource.
func (s *Client) CreateHTTPSource(collectorID int, source HTTPSource) (*HTTPSource, error) {

	request := HTTPSourceRequest{
		Source: source,
	}

	log.Printf("Sumologic API Request: %+v", request)

	body, _ := json.Marshal(request)

	relativeURL, _ := url.Parse(fmt.Sprintf("collectors/%d/sources", collectorID))
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
		var r = new(HTTPSourceRequest)
		err = json.Unmarshal(responseBody, &r)
		if err != nil {
			return nil, err
		}

		return &r.Source, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusBadRequest:
		var e = new(Error)
		return nil, fmt.Errorf("Bad Request. %s", e.Message)
	default:
		return nil, fmt.Errorf("Unknown Response with Sumo Logic: `%d`", resp.StatusCode)
	}
}

// UpdateHTTPSource updates an existing HTTP source.
func (s *Client) UpdateHTTPSource(collectorID int, source HTTPSource, etag string) (*HTTPSource, error) {
	request := HTTPSourceRequest{
		Source: source,
	}

	body, _ := json.Marshal(request)

	relativeURL, _ := url.Parse(fmt.Sprintf("collectors/%d/sources/%d", collectorID, source.ID))
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

	responseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		var r = new(HTTPSourceRequest)
		err = json.Unmarshal(responseBody, &r)
		if err != nil {
			return nil, err
		}

		return &r.Source, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusBadRequest:
		return nil, fmt.Errorf("Bad Request. Please check if a source with this name `%s` already exists", source.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with Sumo Logic: `%d`", resp.StatusCode)
	}
}

// DeleteHTTPSource deletes the source with the specified ID.
func (s *Client) DeleteHTTPSource(collectorID int, id int) error {
	c, _ := url.Parse(fmt.Sprintf("collectors/%d/sources/%d", collectorID, id))
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
		return ErrSourceNotFound
	case http.StatusUnauthorized:
		return ErrClientAuthenticationError
	default:
		return fmt.Errorf("Unknown Response with Sumo Logic: `%d`", resp.StatusCode)
	}
}
