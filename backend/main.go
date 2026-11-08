package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func getDevices(client *http.Client, baseUrl, authStr string) []map[string]any {
	req, err := http.NewRequest("GET", baseUrl+"/devices", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", authStr)
	resp, err := client.Do(req)
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

func getPositions(client *http.Client, baseUrl, authStr string) []map[string]any {
	req, err := http.NewRequest("GET", baseUrl+"/positions", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", authStr)
	resp, err := client.Do(req)
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

func main() {
	fmt.Println("Siya, world!")

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	baseUrl := os.Getenv("TRACCAR_BASE_URL")
	email := os.Getenv("TRACCAR_EMAIL")
	pw := os.Getenv("TRACCAR_PASSWORD")
	authStr := "Basic " + base64.StdEncoding.EncodeToString([]byte(email+":"+pw))

	client := &http.Client{}

	// devices := getDevices(client, baseUrl, authStr)
	// for _, device := range devices {
	// 	fmt.Println("DEVICE")
	// 	for k, v := range device {
	// 		fmt.Println("-", k, ":", v)
	// 	}
	// }

	positions := getPositions(client, baseUrl, authStr)
	for _, position := range positions {
		fmt.Println("POSITION")
		for k, v := range position {
			fmt.Println("-", k, ":", v)
		}
	}
}
