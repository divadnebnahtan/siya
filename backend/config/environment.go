package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const authType = "Basic"

type Environment struct {
	TraccarBaseURL  string
	TraccarEmail    string
	TraccarPassword string
}

func (e Environment) GetAuthString() string {
	return authType + " " + base64.StdEncoding.EncodeToString([]byte(e.TraccarEmail+":"+e.TraccarPassword))
}

func Load() (Environment, error) {
	settings, err := godotenv.Read(".env")
	if err != nil {
		return Environment{}, fmt.Errorf("load .env: %w", err)
	}

	lookup := func(key string) string {
		value, exists := os.LookupEnv(key)
		if exists {
			return value
		}
		return settings[key]
	}
	env := Environment{
		TraccarBaseURL:  lookup("TRACCAR_BASE_URL"),
		TraccarEmail:    lookup("TRACCAR_EMAIL"),
		TraccarPassword: lookup("TRACCAR_PASSWORD"),
	}
	for _, setting := range []struct{ name, value string }{
		{"TRACCAR_BASE_URL", env.TraccarBaseURL},
		{"TRACCAR_EMAIL", env.TraccarEmail},
		{"TRACCAR_PASSWORD", env.TraccarPassword},
	} {
		if setting.value == "" {
			return Environment{}, fmt.Errorf("%s is required", setting.name)
		}
	}
	if strings.Contains(env.TraccarEmail, ":") {
		return Environment{}, fmt.Errorf("TRACCAR_EMAIL cannot contain a colon")
	}
	return env, nil
}
