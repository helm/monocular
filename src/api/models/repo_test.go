package models

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
)

var emptyDB, _ = datastore.NewMockSession(nil, false).DB()

func TestListRepos(t *testing.T) {
	tests := []struct {
		name          string
		expectedRepos []*Repo
	}{
		{"no repos", []*Repo{}},
		{"repos", []*Repo{{Name: "stable", URL: "stable.com", Source: "stable.com"}}},
	}
	for _, tt := range tests {
		db, _ := datastore.NewMockSession(&tt.expectedRepos, false).DB()
		actual, err := ListRepos(db)
		assert.NoErr(t, err)
		assert.Equal(t, actual, tt.expectedRepos, tt.name)
	}
}

func TestGetRepo(t *testing.T) {
	tests := []struct {
		name      string
		expected  *Repo
		expectErr bool
	}{
		{"repo", &Repo{Name: "stable", URL: "stable.com", Source: "stable.com"}, false},
		{"inexistant repo", nil, true},
	}
	for _, tt := range tests {
		db, _ := datastore.NewMockSession(tt.expected, tt.expectErr).DB()
		actual, err := GetRepo(db, "stable")
		assert.Equal(t, err != nil, tt.expectErr, "error")
		assert.Equal(t, actual, tt.expected, "repo")
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
		db, _ := datastore.NewMockSession(nil, tt.expectErr).DB()
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
		db, _ := datastore.NewMockSession(nil, tt.expectErr).DB()
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
		db, _ := datastore.NewMockSession(nil, tt.expectErr).DB()
		assert.Equal(t, DeleteRepo(db, tt.name) != nil, tt.expectErr, tt.name)
	}
}
