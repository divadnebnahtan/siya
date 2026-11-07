package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	baseUrl := os.Getenv("TRACCAR_BASE_URL")
	email := os.Getenv("TRACCAR_EMAIL")
	pw := os.Getenv("TRACCAR_PASSWORD")

	fmt.Println("Go, World!")

	req, err := http.NewRequest("GET", baseUrl+"/devices", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(email+":"+pw)))

	fmt.Println(req.Header.Get("Authorization"))

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status, resp.StatusCode)
	fmt.Println(resp.Proto, resp.ProtoMajor, resp.ProtoMinor)
	for k, v := range resp.Header {
		fmt.Println("-", k, ":", v)
	}
	fmt.Println(resp.Body)
	body, err := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}
