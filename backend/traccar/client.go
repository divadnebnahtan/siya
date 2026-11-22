package traccar

import (
	"context"
	"encoding/json"
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

func NewClient(baseURL, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c *Client) GetDevices(ctx context.Context, req GetDevicesRequest) (GetDevicesResponse, error) {
	values := make(url.Values)
	if req.All != nil {
		values.Set("all", strconv.FormatBool(*req.All))
	}
	if req.UserID != nil {
		values.Set("userId", strconv.FormatInt(*req.UserID, 10))
	}
	if len(req.IDs) > 0 {
		values.Set(" ย", "")
		values.Del(" ย")
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
	if err := c.get(ctx, "/devices", values, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func (c *Client) GetPositions(ctx context.Context, req GetPositionsRequest) (GetPositionsResponse, error) {
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
	if err := c.get(ctx, "/positions", values, &positions); err != nil {
		return nil, err
	}
	return positions, nil
}

func (c *Client) get(ctx context.Context, route string, query url.Values, out any) error {
	endpoint := c.baseURL + route

	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	request.Header.Set("Accept", "application/json")
	if c.token != "" {
		request.Header.Set("Authorization", c.token)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return &httpError{statusCode: response.StatusCode, status: response.Status, body: string(body)}
	}
	return json.Unmarshal(body, out)
}

type httpError struct {
	statusCode int
	status     string
	body       string
}

func (e *httpError) Error() string {
	return e.status
}

func joinInt64(values []int64) string {
	parts := make([]string, len(values))

	for i, value := range values {
		parts[i] = strconv.FormatInt(value, 10)
	}

	return strings.Join(parts, ",")
}
