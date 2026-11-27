package config_test

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/divadnebnahtan/siya/backend/config"
)

func TestLoad(t *testing.T) {
	for _, tc := range []struct {
		name, file, email string
		missing           bool
		cause             error
		wantError         bool
	}{
		{name: "environment takes precedence", file: "TRACCAR_EMAIL=file@example.com", email: "user@example.com"},
		{name: "empty override", file: "TRACCAR_EMAIL=file@example.com", wantError: true},
		{name: "invalid email separator", file: "# test", email: "a:b", wantError: true},
		{name: "invalid dotenv", file: "TRACCAR_EMAIL='unterminated", wantError: true},
		{name: "missing dotenv", missing: true, wantError: true, cause: fs.ErrNotExist},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("TRACCAR_BASE_URL", "https://example.com")
			t.Setenv("TRACCAR_EMAIL", tc.email)
			t.Setenv("TRACCAR_PASSWORD", "password")
			if !tc.missing {
				err := os.WriteFile(".env", []byte(tc.file), 0o600)
				if err != nil {
					t.Fatal(err)
				}
			}
			env, err := config.Load()
			if (err != nil) != tc.wantError || (tc.cause != nil && !errors.Is(err, tc.cause)) {
				t.Fatalf("Load() error = %v, want error: %v", err, tc.wantError)
			}
			if !tc.wantError && env.TraccarEmail != tc.email {
				t.Errorf("Load() email = %q, want %q", env.TraccarEmail, tc.email)
			}
		})
	}
}
