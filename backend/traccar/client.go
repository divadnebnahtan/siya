package traccar

import (
	"encoding/json"
	"io"
	"net/http"
)

type TraccarClient struct {
	httpClient *http.Client
	url        string
	auth       string
}

func NewClient(traccarBaseUrl, traccarAuthString string) *TraccarClient {
	return &TraccarClient{
		httpClient: &http.Client{},
		url:        traccarBaseUrl,
		auth:       traccarAuthString,
	}
}

func (tc *TraccarClient) GetDevices() []map[string]any {
	req, err := http.NewRequest("GET", tc.url+"/devices", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", tc.auth)
	resp, err := tc.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println(resp.Status, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

	var devices [](map[string]any)
	err = json.Unmarshal([]byte(string(body)), &devices)
	if err != nil {
		panic(err)
	}

	return devices
}

func (tc *TraccarClient) GetPositions() []map[string]any {
	req, err := http.NewRequest("GET", tc.url+"/positions", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", tc.auth)
	resp, err := tc.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println(resp.Status, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// fmt.Println(string(body))

	var devices [](map[string]any)
	err = json.Unmarshal([]byte(string(body)), &devices)
	if err != nil {
		panic(err)
	}

	return devices
}
