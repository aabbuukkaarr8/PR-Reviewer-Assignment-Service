package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aabbuukkaarr8/PRService/internal/apiserver"
	pullrequestsHandler "github.com/aabbuukkaarr8/PRService/internal/handler/pullrequests"
	teamHandler "github.com/aabbuukkaarr8/PRService/internal/handler/team"
	usersHandler "github.com/aabbuukkaarr8/PRService/internal/handler/users"
	"github.com/aabbuukkaarr8/PRService/internal/repository/pullrequests"
	"github.com/aabbuukkaarr8/PRService/internal/repository/team"
	"github.com/aabbuukkaarr8/PRService/internal/repository/users"
	pullrequestsService "github.com/aabbuukkaarr8/PRService/internal/service/pullrequests"
	teamService "github.com/aabbuukkaarr8/PRService/internal/service/team"
	usersService "github.com/aabbuukkaarr8/PRService/internal/service/users"
	"github.com/aabbuukkaarr8/PRService/internal/store"
	_ "github.com/lib/pq"
)

var (
	testDB     *sql.DB
	testStore  *store.Store
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://appuser:secret@localhost:5432/PReviewer?sslmode=disable"
	}

	var err error
	testDB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to test database: %v", err))
	}

	if err := testDB.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping test database: %v. Make sure PostgreSQL is running and database exists", err))
	}

	if err := runMigrations(testDB); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	testStore = store.New()
	testStore.SetConn(testDB)

	config := apiserver.NewConfig()
	config.BindAddr = ":0"
	config.LogLevel = "error"

	teamRepo := team.NewRepository(testStore)
	userRepo := users.NewRepository(testStore)
	prRepo := pullrequests.NewRepository(testStore)

	teamSrv := teamService.NewService(teamRepo)
	userSrv := usersService.NewService(userRepo)
	prSrv := pullrequestsService.NewService(prRepo)

	teamHndlr := teamHandler.NewHandler(teamSrv)
	userHndlr := usersHandler.NewHandler(userSrv)
	prHndlr := pullrequestsHandler.NewHandler(prSrv)

	s := apiserver.New(config)
	s.ConfigureRouter(teamHndlr, userHndlr, prHndlr)

	testServer = httptest.NewServer(s.GetRouter())
}

func teardown() {
	if testDB != nil {
		cleanupDatabase(testDB)
		testDB.Close()
	}
	if testServer != nil {
		testServer.Close()
	}
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS teams (
			team_name TEXT PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			team_name TEXT NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
			is_active BOOLEAN NOT NULL DEFAULT TRUE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name)`,
		`CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)`,
		`CREATE TABLE IF NOT EXISTS pullrequests (
			pull_request_id TEXT PRIMARY KEY,
			pull_request_name TEXT NOT NULL,
			author_id TEXT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
			assigned_reviewers TEXT[] DEFAULT '{}',
			"createdAt" TIMESTAMP,
			"mergedAt" TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_pullrequests_author_id ON pullrequests(author_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pullrequests_status ON pullrequests(status)`,
		`CREATE INDEX IF NOT EXISTS idx_pullrequests_assigned_reviewers ON pullrequests USING GIN(assigned_reviewers)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}

func cleanupDatabase(db *sql.DB) {
	tables := []string{"pullrequests", "users", "teams"}
	for _, table := range tables {
		db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
	}
}

func TestE2E_CreateTeamAndGetTeam(t *testing.T) {
	cleanupDatabase(testDB)

	client := &http.Client{Timeout: 5 * time.Second}

	teamData := map[string]interface{}{
		"team_name": "backend",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
		},
	}

	body, _ := json.Marshal(teamData)
	req, _ := http.NewRequest("POST", testServer.URL+"/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	req2, _ := http.NewRequest("GET", testServer.URL+"/team/get?team_name=backend", nil)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp2.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["team_name"] != "backend" {
		t.Errorf("Expected team_name 'backend', got %v", result["team_name"])
	}
}

func TestE2E_CreatePR(t *testing.T) {
	cleanupDatabase(testDB)

	teamData := map[string]interface{}{
		"team_name": "backend",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
			{"user_id": "u3", "username": "Charlie", "is_active": true},
		},
	}

	client := &http.Client{Timeout: 5 * time.Second}

	body, _ := json.Marshal(teamData)
	req, _ := http.NewRequest("POST", testServer.URL+"/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create team: %v", err)
	}
	resp.Body.Close()

	prData := map[string]interface{}{
		"pull_request_id":   "pr-1001",
		"pull_request_name": "Add feature",
		"author_id":         "u1",
	}

	body, _ = json.Marshal(prData)
	req, _ = http.NewRequest("POST", testServer.URL+"/pullRequest/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create PR: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	pr, ok := result["pr"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'pr' in response")
	}

	if pr["status"] != "OPEN" {
		t.Errorf("Expected status 'OPEN', got %v", pr["status"])
	}

	reviewers, ok := pr["assigned_reviewers"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'assigned_reviewers' in response")
	}

	if len(reviewers) == 0 {
		t.Error("Expected at least one reviewer")
	}
}

func TestE2E_MergePR(t *testing.T) {
	cleanupDatabase(testDB)

	client := &http.Client{Timeout: 5 * time.Second}

	teamData := map[string]interface{}{
		"team_name": "backend",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
		},
	}

	body, _ := json.Marshal(teamData)
	req, _ := http.NewRequest("POST", testServer.URL+"/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	resp.Body.Close()

	prData := map[string]interface{}{
		"pull_request_id":   "pr-1001",
		"pull_request_name": "Add feature",
		"author_id":         "u1",
	}

	body, _ = json.Marshal(prData)
	req, _ = http.NewRequest("POST", testServer.URL+"/pullRequest/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	resp.Body.Close()

	mergeData := map[string]interface{}{
		"pull_request_id": "pr-1001",
	}

	body, _ = json.Marshal(mergeData)
	req, _ = http.NewRequest("POST", testServer.URL+"/pullRequest/merge", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to merge PR: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	pr, ok := result["pr"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'pr' in response")
	}

	if pr["status"] != "MERGED" {
		t.Errorf("Expected status 'MERGED', got %v", pr["status"])
	}
}

func TestE2E_SetUserActive(t *testing.T) {
	cleanupDatabase(testDB)

	client := &http.Client{Timeout: 5 * time.Second}

	teamData := map[string]interface{}{
		"team_name": "backend",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
		},
	}

	body, _ := json.Marshal(teamData)
	req, _ := http.NewRequest("POST", testServer.URL+"/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	resp.Body.Close()

	userData := map[string]interface{}{
		"user_id":   "u1",
		"is_active": false,
	}

	body, _ = json.Marshal(userData)
	req, _ = http.NewRequest("POST", testServer.URL+"/users/setIsActive", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to set user active: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	user, ok := result["user"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'user' in response")
	}

	if user["is_active"] != false {
		t.Errorf("Expected is_active false, got %v", user["is_active"])
	}
}

func TestE2E_GetUserReviews(t *testing.T) {
	cleanupDatabase(testDB)

	client := &http.Client{Timeout: 5 * time.Second}

	teamData := map[string]interface{}{
		"team_name": "backend",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
		},
	}

	body, _ := json.Marshal(teamData)
	req, _ := http.NewRequest("POST", testServer.URL+"/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	resp.Body.Close()

	prData := map[string]interface{}{
		"pull_request_id":   "pr-1001",
		"pull_request_name": "Add feature",
		"author_id":         "u1",
	}

	body, _ = json.Marshal(prData)
	req, _ = http.NewRequest("POST", testServer.URL+"/pullRequest/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	resp.Body.Close()

	req, _ = http.NewRequest("GET", testServer.URL+"/users/getReview?user_id=u2", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get user reviews: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["user_id"] != "u2" {
		t.Errorf("Expected user_id 'u2', got %v", result["user_id"])
	}
}

func TestE2E_ReassignReviewer(t *testing.T) {
	cleanupDatabase(testDB)

	client := &http.Client{Timeout: 5 * time.Second}

	teamData := map[string]interface{}{
		"team_name": "backend",
		"members": []map[string]interface{}{
			{"user_id": "u1", "username": "Alice", "is_active": true},
			{"user_id": "u2", "username": "Bob", "is_active": true},
			{"user_id": "u3", "username": "Charlie", "is_active": true},
			{"user_id": "u4", "username": "David", "is_active": true},
		},
	}

	body, _ := json.Marshal(teamData)
	req, _ := http.NewRequest("POST", testServer.URL+"/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := client.Do(req)
	resp.Body.Close()

	prData := map[string]interface{}{
		"pull_request_id":   "pr-1001",
		"pull_request_name": "Add feature",
		"author_id":         "u1",
	}

	body, _ = json.Marshal(prData)
	req, _ = http.NewRequest("POST", testServer.URL+"/pullRequest/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	resp.Body.Close()

	reassignData := map[string]interface{}{
		"pull_request_id": "pr-1001",
		"old_reviewer_id": "u2",
	}

	body, _ = json.Marshal(reassignData)
	req, _ = http.NewRequest("POST", testServer.URL+"/pullRequest/reassign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to reassign reviewer: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		t.Errorf("Expected status 200, got %d. Body: %+v", resp.StatusCode, errorBody)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, ok := result["replaced_by"]; !ok {
		t.Error("Expected 'replaced_by' in response")
	}
}
