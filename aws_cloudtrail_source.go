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

type Error struct {
	Status  int    `json:"status"`
	ID      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AWSCloudTrailSource is a necessary wrapper for source API calls.
type AWSCloudTrailSourceRequest struct {
	Source AWSCloudTrailSource `json:"source"`
}

// AWSCloudTrailSource can various types of sources including Cloudtrail and S3.
type AWSCloudTrailSource struct {
	ID                 int                    `json:"id,omitempty"`
	Name               string                 `json:"name"`
	CollectorID        int                    `json:"CollectorId,omitempty"`
	Description        string                 `json:"description,omitempty"`
	Category           string                 `json:"category,omitempty"`
	TimeZone           string                 `json:"timezone,omitempty"`
	SourceType         string                 `json:"sourceType,omitempty"`
	ContentType        string                 `json:"contentType,omitempty"`
	ScanInterval       int                    `json:"scanInterval,omitempty"`
	Paused             bool                   `json:"paused"`
	CutoffRelativeTime string                 `json:"cutoffRelativeTime"`
	ThirdPartyRef      AWSBucketThirdPartyRef `json:"thirdPartyRef,omitempty"`
}

type AWSBucketThirdPartyRef struct {
	Resources []AWSBucketResource `json:"resources,omitempty"`
}

// AWSBucketThirdPartyRef contains AWS configurartion including auth.
type AWSBucketResource struct {
	ServiceType    string                  `json:"serviceType"`
	Path           AWSBucketPath           `json:"path"`
	Authentication AWSBucketAuthentication `json:"authentication"`
}

// AWSBucketPath contains AWS S3 Bucket configuration.
type AWSBucketPath struct {
	Type           string `json:"type"`
	BucketName     string `json:"bucketName"`
	PathExpression string `json:"pathExpression"`
}

// AWSBucketAuthentication contains AWS authentication configurartion.
type AWSBucketAuthentication struct {
	Type    string `json:"type"`
	RoleARN string `json:"roleARN"`
}

// ErrSourceNotFound is returned when a source doesn't exist on a Read or Delete.
// It's useful for ignoring errors (e.g. delete if exists).
var ErrSourceNotFound = errors.New("Source not found")

// ErrAwsAuthenticationError is returned for authentication errors with AWS.
// Due to IAM's eventual consistency, it may be useful to retry.
var ErrAwsAuthenticationError = errors.New("Authentication Error with Sumo Logic")

// GetAWSCloudTrailSource gets the source with the specified ID.
func (s *Client) GetAWSCloudTrailSource(collectorID int, id int) (*AWSCloudTrailSource, string, error) {

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
		var r = new(AWSCloudTrailSourceRequest)
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

// CreateAWSCloudTrailSource creates a new AWSCloudTrailSource.
func (s *Client) CreateAWSCloudTrailSource(collectorID int, source AWSCloudTrailSource) (*AWSCloudTrailSource, error) {

	request := AWSCloudTrailSourceRequest{
		Source: source,
	}

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
		var r = new(AWSCloudTrailSourceRequest)
		err = json.Unmarshal(responseBody, &r)
		if err != nil {
			return nil, err
		}

		return &r.Source, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusBadRequest:
		var e = new(Error)
		err = json.Unmarshal(responseBody, &e)
		if err != nil {
			return nil, fmt.Errorf("Bad Request. Please check if a source with this name `%s` already exists", source.Name)
		}
		if e.Message == "Cannot authenticate with AWS." ||
			e.Message == "Invalid IAM role: 'errorCode=AccessDenied'." {
			return nil, ErrAwsAuthenticationError
		}
		return nil, fmt.Errorf("Bad Request. %s", e.Message)
	default:
		return nil, fmt.Errorf("Unknown Response with Sumo Logic: `%d`", resp.StatusCode)
	}
}

// UpdateAWSCloudTrailSource updates an existing AWS Bucket source.
func (s *Client) UpdateAWSCloudTrailSource(collectorID int, source AWSCloudTrailSource, etag string) (*AWSCloudTrailSource, error) {
	request := AWSCloudTrailSourceRequest{
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
		var r = new(AWSCloudTrailSourceRequest)
		err = json.Unmarshal(responseBody, &r)
		if err != nil {
			return nil, err
		}

		return &r.Source, nil
	case http.StatusUnauthorized:
		return nil, ErrClientAuthenticationError
	case http.StatusBadRequest:
		var e = new(Error)
		err = json.Unmarshal(responseBody, &e)
		if e.Message == "Cannot authenticate with AWS." ||
			e.Message == "Invalid IAM role: 'errorCode=AccessDenied'." {
			return nil, ErrAwsAuthenticationError
		}
		return nil, fmt.Errorf("Bad Request. Please check if a source with this name `%s` already exists", source.Name)
	default:
		return nil, fmt.Errorf("Unknown Response with Sumo Logic: `%d`", resp.StatusCode)
	}
}

// DeleteAWSCloudTrailSource deletes the source with the specified ID.
func (s *Client) DeleteAWSCloudTrailSource(collectorID int, id int) error {
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
