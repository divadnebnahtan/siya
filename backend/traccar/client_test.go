package traccar_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/divadnebnahtan/siya/backend/traccar"
)

func TestNewClient(t *testing.T) {
	for _, tc := range []struct {
		name, endpoint, auth string
		wantError            bool
	}{
		{name: "valid", endpoint: "https://example.com/api", auth: "Basic abc"},
		{name: "relative", endpoint: "/api", wantError: true},
		{name: "malformed", endpoint: "http://[::1", wantError: true},
		{name: "unsupported scheme", endpoint: "ftp://example.com", wantError: true},
		{name: "userinfo", endpoint: "https://secret@example.com", wantError: true},
		{name: "query", endpoint: "https://example.com?q=value", wantError: true},
		{name: "invalid header", endpoint: "https://example.com", auth: "bad\r\nvalue", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, err := traccar.NewClient(tc.endpoint, tc.auth)
			if (err != nil) != tc.wantError || (client == nil) != tc.wantError {
				t.Fatalf("traccar.NewClient() = %v, %v, want error: %v", client != nil, err, tc.wantError)
			}
		})
	}
}

func TestGetDevices(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		wantQuery  string
		request    traccar.GetDevicesRequest
	}{
		{name: "success", body: `[{"id":7}]`, status: 200},
		{name: "filters", body: `[{"id":7}]`, status: 200, wantQuery: "id=7%2C9&keyword=a+b", request: traccar.GetDevicesRequest{IDs: []int64{7, 9}, Keyword: "a b"}},
		{name: "status error", body: "denied", status: 401},
		{name: "invalid JSON", body: "[", status: 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.URL.Path != "/api/devices" || r.URL.RawQuery != tc.wantQuery ||
					r.Header.Get("Accept") != "application/json" || r.Header.Get("Authorization") != "Basic abc" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				w.WriteHeader(tc.status)

				_, err := fmt.Fprint(w, tc.body)
				if err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()

			client, err := traccar.NewClient(server.URL+"/api", "Basic abc")
			if err != nil {
				t.Fatal(err)
			}
			devices, err := client.GetDevices(tc.request)
			switch tc.name {
			case "status error":
				var statusErr *traccar.HTTPError
				if !errors.As(err, &statusErr) || statusErr.StatusCode != tc.status || statusErr.Body != tc.body {
					t.Fatalf("expected HTTP status error, got %v", err)
				}
			case "invalid JSON":
				var syntaxErr *json.SyntaxError
				if !errors.As(err, &syntaxErr) {
					t.Fatalf("expected JSON syntax error, got %v", err)
				}
			default:
				if err != nil || len(devices) != 1 || devices[0].ID != 7 {
					t.Fatalf("GetDevices() = %v, %v", devices, err)
				}
			}
		})
	}
}

func TestGetPositionsValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		client *traccar.Client
	}{
		{name: "nil client"},
		{name: "zero client", client: &traccar.Client{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			positions, err := tc.client.GetPositions(traccar.GetPositionsRequest{})
			if err == nil || positions != nil {
				t.Fatalf("GetPositions() = %v, %v", positions, err)
			}
		})
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestResponseReadError(t *testing.T) {
	readErr := errors.New("broken reader")
	client, err := traccar.NewClient("https://example.com", "")
	if err != nil {
		t.Fatal(err)
	}
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	http.DefaultTransport = transportFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(failingReader{readErr})}, nil
	})
	_, err = client.GetPositions(traccar.GetPositionsRequest{})
	if !errors.Is(err, readErr) {
		t.Fatalf("expected read error, got %v", err)
	}
}

func TestGetPositionsConvertsSpeed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"id":1,"speed":10},{"id":2,"speed":0}]`)
	}))
	defer server.Close()

	client, err := traccar.NewClient(server.URL, "")
	if err != nil {
		t.Fatal(err)
	}
	positions, err := client.GetPositions(traccar.GetPositionsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(positions) != 2 || math.Abs(positions[0].Speed-18.52) > 1e-9 || positions[1].Speed != 0 {
		t.Fatalf("unexpected positions: %+v", positions)
	}
}
