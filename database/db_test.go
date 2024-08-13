package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDatabase(t *testing.T) *Database {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return &Database{
		Db: db,
	}
}

func TestNewDatabase(t *testing.T) {
	db := setupTestDatabase(t)
	assert.NotNil(t, db)
}

func TestNewUser(t *testing.T) {
	db := setupTestDatabase(t)

	user := &User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Age:      30,
	}

	createdUser, err := db.NewUser(user)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "testuser", createdUser.Username)
}

func TestGetUser(t *testing.T) {
	db := setupTestDatabase(t)

	user := &User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Age:      30,
	}

	createdUser, err := db.NewUser(user)
	assert.NoError(t, err)

	retrievedUser, err := db.GetUser(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, retrievedUser.ID)
	assert.Equal(t, "testuser", retrievedUser.Username)
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDatabase(t)

	user := &User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Age:      30,
	}

	createdUser, err := db.NewUser(user)
	assert.NoError(t, err)

	updateData := &User{
		Username: "updateduser",
		Email:    "updateduser@example.com",
		Age:      31,
	}

	err = db.UpdateUser(createdUser.ID, updateData)
	assert.NoError(t, err)

	updatedUser, err := db.GetUser(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUser.Username)
	assert.Equal(t, "updateduser@example.com", updatedUser.Email)
	assert.Equal(t, 31, updatedUser.Age)
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDatabase(t)

	user := &User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Age:      30,
	}

	createdUser, err := db.NewUser(user)
	assert.NoError(t, err)

	err = db.DeleteUser(createdUser.ID)
	assert.NoError(t, err)

	_, err = db.GetUser(createdUser.ID)
	assert.Error(t, err)
}
