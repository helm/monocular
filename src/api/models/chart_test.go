package models

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
)

func TestCreateChart(t *testing.T) {
	chart := &Chart{Name: "drupal", Repo: "stable"}
	tests := []struct {
		name      string
		expectErr bool
	}{
		{"create chart", false},
		{"error", true},
	}
	for _, tt := range tests {
		db, _ := datastore.NewMockSession(nil, tt.expectErr).DB()
		assert.Equal(t, CreateChart(db, chart) != nil, tt.expectErr, tt.name)
	}
}

func TestGetChartByName(t *testing.T) {
	tests := []struct {
		name      string
		expected  *Chart
		expectErr bool
	}{
		{"chart", &Chart{Name: "drupal", Repo: "stable"}, false},
		{"inexistant chart", nil, true},
	}
	for _, tt := range tests {
		db, _ := datastore.NewMockSession(tt.expected, tt.expectErr).DB()
		actual, err := GetChartByName(db, "stable", "drupal")
		assert.Equal(t, err != nil, tt.expectErr, "error")
		assert.Equal(t, actual, tt.expected, "repo")
	}
}
