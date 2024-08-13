package api

import (
	"atmail-demo/config"
	"atmail-demo/database"
	"bytes"
	"encoding/json"
	"github.com/nokusukun/faust"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestDatabase(t *testing.T) *database.Database {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&database.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return &database.Database{
		Db: db,
	}
}

func setupServer(db *database.Database) *faust.API {
	server := faust.New()
	UsersEndpoint(server, db)
	return server
}

func TestUserPayloadValidator(t *testing.T) {
	tests := []struct {
		name    string
		payload UserPayload
		wantErr bool
	}{
		{"Valid payload", UserPayload{Username: "testuser", Email: "test@example.com", Age: 25}, false},
		{"Missing username", UserPayload{Email: "test@example.com", Age: 25}, true},
		{"Invalid email", UserPayload{Username: "testuser", Email: "test", Age: 25}, true},
		{"Invalid age", UserPayload{Username: "testuser", Email: "test@example.com", Age: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserPayloadValidator(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserPayloadValidator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserEndpoint(t *testing.T) {
	db := setupTestDatabase(t)
	server := setupServer(db)

	// Insert a user to retrieve
	user := &database.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Age:      30,
	}
	db.NewUser(user)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	req.Header.Set("Authorization", "Basic "+config.USERNAME+":"+config.PASSWORD)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseUser database.User
	err := json.NewDecoder(rr.Body).Decode(&responseUser)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, responseUser.Username)
}

func TestSomething(t *testing.T) {

	payload := UserPayload{
		Username: "newuser",
		Email:    "newuser@example.com",
		Age:      22,
	}
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/users/", bytes.NewBuffer(body))
	assert.NoError(t, err)

	newPayload := UserPayload{}
	json.NewDecoder(req.Body).Decode(&newPayload)
	assert.Equal(t, payload, newPayload)
}

func TestCreateUserEndpoint(t *testing.T) {
	db := setupTestDatabase(t)
	server := setupServer(db)

	payload := UserPayload{
		Username: "newuser",
		Email:    "newuser@example.com",
		Age:      22,
	}
	body, err := json.Marshal(payload)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/users/", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+config.USERNAME+":"+config.PASSWORD)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response["ok"].(bool))
}

func TestUpdateUserEndpoint(t *testing.T) {
	db := setupTestDatabase(t)
	server := setupServer(db)

	// Insert a user to update
	user := &database.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Age:      30,
	}
	db.NewUser(user)

	updatePayload := UserPayload{
		Username: "updateduser",
		Email:    "updateduser@example.com",
		Age:      31,
	}
	body, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+config.USERNAME+":"+config.PASSWORD)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteUserEndpoint(t *testing.T) {
	db := setupTestDatabase(t)
	server := setupServer(db)

	// Insert a user to delete
	user := &database.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Age:      30,
	}
	db.NewUser(user)

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	req.Header.Set("Authorization", "Basic "+config.USERNAME+":"+config.PASSWORD)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
