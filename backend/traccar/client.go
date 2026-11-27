package traccar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string) (*Client, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		// The parse error includes the input, which may contain credentials.
		return nil, errors.New("invalid Traccar base URL")
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Hostname() == "" ||
		parsed.User != nil || parsed.Opaque != "" || parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" {
		return nil, errors.New("base URL for Traccar must be an absolute HTTP(S) URL without credentials, query, or fragment")
	}
	for i := range len(token) {
		if token[i] < ' ' && token[i] != '\t' || token[i] == 0x7f {
			return nil, errors.New("invalid Traccar authorization header")
		}
	}

	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{},
	}, nil
}

func (c *Client) GetDevices(req GetDevicesRequest) (GetDevicesResponse, error) {
	values := make(url.Values)
	if req.All != nil {
		values.Set("all", strconv.FormatBool(*req.All))
	}
	if req.UserID != nil {
		values.Set("userId", strconv.FormatInt(*req.UserID, 10))
	}
	if len(req.IDs) > 0 {
		values.Set("id", joinInt64(req.IDs))
	}
	if len(req.UniqueIDs) > 0 {
		values.Set("uniqueId", strings.Join(req.UniqueIDs, ","))
	}
	if req.ExcludeAttributes != nil {
		values.Set("excludeAttributes", strconv.FormatBool(*req.ExcludeAttributes))
	}
	if req.Limit != nil {
		values.Set("limit", strconv.Itoa(*req.Limit))
	}
	if req.Offset != nil {
		values.Set("offset", strconv.Itoa(*req.Offset))
	}
	if req.Keyword != "" {
		values.Set("keyword", req.Keyword)
	}

	var devices GetDevicesResponse
	if err := c.get("/devices", values, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func (c *Client) GetPositions(req GetPositionsRequest) (GetPositionsResponse, error) {
	values := make(url.Values)
	if req.DeviceID != nil {
		values.Set("deviceId", strconv.FormatInt(*req.DeviceID, 10))
	}
	if req.From != nil {
		values.Set("from", req.From.Format(time.RFC3339Nano))
	}
	if req.To != nil {
		values.Set("to", req.To.Format(time.RFC3339Nano))
	}
	if len(req.IDs) > 0 {
		values.Set("id", joinInt64(req.IDs))
	}

	var positions GetPositionsResponse
	if err := c.get("/positions", values, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}

func (c *Client) get(route string, query url.Values, out any) error {
	if c == nil || c.httpClient == nil || c.baseURL == "" {
		return errors.New("uninitialized Traccar client")
	}

	endpoint := c.baseURL + route
	if encoded := query.Encode(); encoded != "" {
		endpoint += "?" + encoded
	}
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build Traccar GET %s: %w", route, err)
	}
	request.Header.Set("Accept", "application/json")
	if c.token != "" {
		request.Header.Set("Authorization", c.token)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send Traccar GET %s: %w", route, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read Traccar GET %s response: %w", route, err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		statusErr := &HTTPError{StatusCode: response.StatusCode, Status: response.Status, Body: string(body)}
		return fmt.Errorf("unexpected status for Traccar GET %s: %w", route, statusErr)
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode Traccar GET %s response: %w", route, err)
	}
	return nil
}

// HTTPError describes a non-success response from Traccar.
type HTTPError struct {
	StatusCode int
	Status     string
	Body       string
}

func (e *HTTPError) Error() string {
	return e.Status
}

func joinInt64(values []int64) string {
	parts := make([]string, len(values))

	for i, value := range values {
		parts[i] = strconv.FormatInt(value, 10)
	}

	return strings.Join(parts, ",")
}
