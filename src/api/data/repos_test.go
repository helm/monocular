package data

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
)

func TestRepo_ModelId(t *testing.T) {
	tests := []struct {
		name string
		r    *Repo
		want string
	}{
		{"stable repo id", &Repo{Name: pointerto.String("stable")}, "stable"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.r.ModelId(), tt.want, tt.name)
		})
	}
}

func TestRepo_SetModelId(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		r    *Repo
		args args
	}{
		{"stable repo id", &Repo{}, args{"stable"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.SetModelId(tt.args.name)
			assert.Equal(t, *tt.r.Name, tt.args.name, tt.name)
		})
	}
}
