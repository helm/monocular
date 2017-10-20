package models

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
)

func TestCreateUser(t *testing.T) {
	user := &User{Name: "Rick Sanchez", Email: "rick@sanchez.com"}
	tests := []struct {
		name      string
		expectErr bool
	}{
		{"create user", false},
		{"error", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := datastore.NewMockSession(nil, tt.expectErr).DB()
			err := CreateUser(db, user)
			assert.Equal(t, err != nil, tt.expectErr, tt.name)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	tests := []struct {
		name      string
		expected  *User
		expectErr bool
	}{
		{"repo", &User{Name: "Rick Sanchez", Email: "rick@sanchez.com"}, false},
		{"inexistant repo", nil, true},
	}
	for _, tt := range tests {
		db, _ := datastore.NewMockSession(tt.expected, tt.expectErr).DB()
		actual, err := GetUserByEmail(db, "rick@sanchez.com")
		assert.Equal(t, err != nil, tt.expectErr, "error")
		assert.Equal(t, actual, tt.expected, "repo")
	}
}
