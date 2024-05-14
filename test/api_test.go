package test

import (
	"bash_command_api/config"
	"bash_command_api/internal/api"
	"bash_command_api/internal/db"
	"bash_command_api/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateCommand(t *testing.T) {
	cfg, _ := config.LoadConfig()
	db.InitDB(cfg)

	handler := http.HandlerFunc(api.CreateCommand)
	server := httptest.NewServer(handler)
	defer server.Close()

	jsonData := `{"command": "echo Hello World"}`
	req, err := http.NewRequest("POST", server.URL, strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %d, but got %d", http.StatusCreated, resp.StatusCode)
	}

	var command models.Command
	if err := json.NewDecoder(resp.Body).Decode(&command); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}
	if command.Command != "echo Hello World" {
		t.Fatalf("Expected command 'echo Hello World', but got '%s'", command.Command)
	}
}

func TestGetCommands(t *testing.T) {
	cfg, _ := config.LoadConfig()
	db.InitDB(cfg)

	handler := http.HandlerFunc(api.GetCommands)
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
}
