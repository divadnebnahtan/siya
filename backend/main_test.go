package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/divadnebnahtan/siya/backend/state"
)

func TestRun(t *testing.T) {
	for _, tc := range []struct {
		name, positions string
		wantError       bool
	}{
		{name: "success", positions: `[
			{"id":9,"deviceId":7,"fixTime":"2026-01-02T00:00:00Z","speed":10,"attributes":{"motion":false}},
			{"id":8,"deviceId":7,"fixTime":"2026-01-01T00:00:00Z","speed":5},
			{"id":10,"deviceId":99}
		]`},
		{name: "invalid positions", positions: "[", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			err := os.WriteFile(".env", nil, 0o600)
			if err != nil {
				t.Fatal(err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.RawQuery != "" {
					t.Errorf("unexpected query: %q", r.URL.RawQuery)
				}
				var body string
				switch r.URL.Path {
				case "/api/devices":
					body = `[{"id":8,"name":"No fix"},{"id":7,"name":"test"}]`
				case "/api/positions":
					body = tc.positions
				default:
					t.Errorf("unexpected route: %s", r.URL.Path)
					http.NotFound(w, r)
					return
				}

				_, err := fmt.Fprint(w, body)
				if err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()
			t.Setenv("TRACCAR_BASE_URL", server.URL+"/api")
			t.Setenv("TRACCAR_EMAIL", "user@example.com")
			t.Setenv("TRACCAR_PASSWORD", "password")

			path := filepath.Join(t.TempDir(), "stdout")
			out, err := os.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			runErr := run(out)
			err = out.Close()
			if err != nil {
				t.Fatal(err)
			}
			if (runErr != nil) != tc.wantError {
				t.Fatalf("run() error = %v, want error: %v", runErr, tc.wantError)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			if tc.wantError {
				if len(data) != 0 {
					t.Fatalf("unexpected output on error: %s", data)
				}
				return
			}
			var locations []state.Location

			err = json.Unmarshal(data, &locations)
			if err != nil {
				t.Fatal(err)
			}
			if len(locations) != 2 || locations[0].ID != 7 || locations[0].Name != "test" ||
				locations[0].Source != "traccar" || locations[1].ID != 8 || locations[1].Position != nil {
				t.Fatalf("unexpected snapshot: %s", data)
			}
			position := locations[0].Position
			if position == nil || position.Speed != 18.52 || position.Motion == nil ||
				*position.Motion || position.ReceivedAt.IsZero() {
				t.Fatalf("unexpected position: %s", data)
			}
		})
	}
}
