package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
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
	ChartVersions []chartVersion
}

type chartVersion struct {
	ID         bson.ObjectId `bson:"_id"`
	Version    string
	AppVersion string
	Created    time.Time
	Digest     string
	Readme     string
	URLs       []string
}
