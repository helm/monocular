package models

import (
	"testing"

	"github.com/arschles/assert"
)

func TestListRepos(t *testing.T) {
	tests := []struct {
		name          string
		expectedRepos []*Repo
	}{
		{"repos", OfficialRepos},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewMockSession(MockDBConfig{}).DB()
			actual, err := ListRepos(db)
			assert.NoErr(t, err)
			assert.Equal(t, actual, tt.expectedRepos, tt.name)
		})
	}
}

func TestGetRepo(t *testing.T) {
	tests := []struct {
		name      string
		expected  *Repo
		expectErr bool
	}{
		{"repo", OfficialRepos[0], false},
		{"inexistant repo", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewMockSession(MockDBConfig{WantErr: tt.expectErr}).DB()
			actual, err := GetRepo(db, "stable")
			assert.Equal(t, err != nil, tt.expectErr, "error")
			assert.Equal(t, actual, tt.expected, "repo")
		})
	}
}

func TestCreateRepos(t *testing.T) {
	repos := []*Repo{
		{Name: "stable", URL: "stable.com", Source: "stable.com"},
		{Name: "incubator", URL: "incubator.com", Source: "incubator.com"},
	}
	tests := []struct {
		name      string
		expectErr bool
	}{
		{"create repos", false},
		{"error", true},
	}
	for _, tt := range tests {
		db, _ := NewMockSession(MockDBConfig{WantErr: tt.expectErr}).DB()
		assert.Equal(t, CreateRepos(db, repos) != nil, tt.expectErr, tt.name)
	}
}

func TestCreateRepo(t *testing.T) {
	repo := &Repo{Name: "stable", URL: "stable.com", Source: "stable.com"}
	tests := []struct {
		name      string
		expectErr bool
	}{
		{"create repos", false},
		{"error", true},
	}
	for _, tt := range tests {
		db, _ := NewMockSession(MockDBConfig{WantErr: tt.expectErr}).DB()
		assert.Equal(t, CreateRepo(db, repo) != nil, tt.expectErr, tt.name)
	}
}

func TestDeleteRepo(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
	}{
		{"stable", false},
		{"inexistant", true},
	}
	for _, tt := range tests {
		db, _ := NewMockSession(MockDBConfig{WantErr: tt.expectErr}).DB()
		assert.Equal(t, DeleteRepo(db, tt.name) != nil, tt.expectErr, tt.name)
	}
}
