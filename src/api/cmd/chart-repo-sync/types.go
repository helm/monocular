package main

import (
	"time"
)

type repo struct {
	Name string
	URL  string
}

type maintainer struct {
	Name  string
	Email string
}

type chart struct {
	ID            string `bson:"_id"`
	Name          string
	Repo          repo
	Description   string
	Home          string
	Keywords      []string
	Maintainers   []maintainer
	Sources       []string
	Icon          string
	ChartVersions []chartVersion
}

type chartVersion struct {
	Version    string
	AppVersion string
	Created    time.Time
	Digest     string
	URLs       []string
}

type chartFiles struct {
	ID     string `bson:"_id"`
	Readme string
	Values string
}
