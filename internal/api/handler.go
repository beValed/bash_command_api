package api

import (
	"bash_command_api/internal/db"
	"bash_command_api/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var commandExecutions = make(map[uint]*exec.Cmd)
var mu sync.Mutex

func CreateCommand(w http.ResponseWriter, r *http.Request) {
	var command models.Command
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	command.Status = "Running"
	command.CreatedAt = time.Now()

	query := `INSERT INTO commands (command, status, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := db.DB.QueryRow(query, command.Command, command.Status, command.CreatedAt).Scan(&command.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go func(cmd models.Command) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		execCmd := exec.CommandContext(ctx, "bash", "-c", cmd.Command)
		mu.Lock()
		commandExecutions[cmd.ID] = execCmd
		mu.Unlock()
		out, err := execCmd.CombinedOutput()
		mu.Lock()
		delete(commandExecutions, cmd.ID)
		mu.Unlock()
		if err != nil {
			cmd.Status = "Failed"
			cmd.Output = err.Error()
		} else {
			cmd.Status = "Completed"
			cmd.Output = string(out)
		}
		updateQuery := `UPDATE commands SET status = $1, output = $2 WHERE id = $3`
		_, err = db.DB.Exec(updateQuery, cmd.Status, cmd.Output, cmd.ID)
		if err != nil {
			log.Printf("Failed to update command: %v", err)
		}
	}(command)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(command)
}

func GetCommands(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, command, status, output, created_at FROM commands")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commands []models.Command
	for rows.Next() {
		var command models.Command
		if err := rows.Scan(&command.ID, &command.Command, &command.Status, &command.Output, &command.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commands = append(commands, command)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commands)
}

func GetCommand(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var command models.Command
	query := `SELECT id, command, status, output, created_at FROM commands WHERE id = $1`
	err = db.DB.QueryRow(query, id).Scan(&command.ID, &command.Command, &command.Status, &command.Output, &command.CreatedAt)
	if err != nil {
		http.Error(w, "Command not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(command)
}

func StopCommand(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var command models.Command
	query := `SELECT id, command, status, output, created_at FROM commands WHERE id = $1`
	err = db.DB.QueryRow(query, id).Scan(&command.ID, &command.Command, &command.Status, &command.Output, &command.CreatedAt)
	if err != nil {
		http.Error(w, "Command not found", http.StatusNotFound)
		return
	}

	mu.Lock()
	execCmd, exists := commandExecutions[command.ID]
	mu.Unlock()

	if exists {
		execCmd.Process.Kill()
		command.Status = "Stopped"
		updateQuery := `UPDATE commands SET status = $1 WHERE id = $2`
		_, err = db.DB.Exec(updateQuery, command.Status, command.ID)
		if err != nil {
			http.Error(w, "Failed to stop command", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "Command stopped"})
	} else {
		http.Error(w, "Command already completed or not found", http.StatusBadRequest)
	}
}
