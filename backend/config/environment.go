package config

import (
	"encoding/base64"
	"os"

	"github.com/joho/godotenv"
)

const AUTH_TYPE string = "Basic"

type Environment struct {
	TraccarBaseUrl  string
	TraccarEmail    string
	TraccarPassword string
}

func (e Environment) GetAuthString() string {
	return AUTH_TYPE + " " + base64.StdEncoding.EncodeToString([]byte(e.TraccarEmail+":"+e.TraccarPassword))
}

func Load() Environment {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	environment := Environment{
		TraccarBaseUrl:  os.Getenv("TRACCAR_BASE_URL"),
		TraccarEmail:    os.Getenv("TRACCAR_EMAIL"),
		TraccarPassword: os.Getenv("TRACCAR_PASSWORD"),
	}
	return environment

}
