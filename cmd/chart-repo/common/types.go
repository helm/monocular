/*
Copyright (c) 2019 The Helm Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"net/http"
	"time"
)

//Repo holds information to identify a repository
type Repo struct {
	Name                string
	URL                 string
	AuthorizationHeader string `bson:"-"`
}

//Maintainer describes the maintainer of a Chart
type Maintainer struct {
	Name  string
	Email string
}

//Chart holds full descriptor of a Helm chart
type Chart struct {
	ID            string `bson:"_id"`
	Name          string
	Repo          Repo
	Description   string
	Home          string
	Keywords      []string
	Maintainers   []Maintainer
	Sources       []string
	Icon          string
	ChartVersions []ChartVersion
}

//ChartVersion holds version information on a Chart
type ChartVersion struct {
	Version    string
	AppVersion string
	Created    time.Time
	Digest     string
	URLs       []string
}

//ChartFiles describes the chart values, readme, schema and digest components of a chart
type ChartFiles struct {
	ID     string `bson:"_id"`
	Readme string
	Values string
	Schema string
	Repo   Repo
	Digest string
}

//RepoCheck describes the state of a repository in terms its current checksum and last update time.
//It is used to determine whether or not to re-sync a respository.
type RepoCheck struct {
	ID         string    `bson:"_id"`
	LastUpdate time.Time `bson:"last_update"`
	Checksum   string    `bson:"checksum"`
}

//ImportChartFilesJob contains the information needed by an
//ImportWorker when import a chart from a repository
type ImportChartFilesJob struct {
	Name         string
	Repo         Repo
	ChartVersion ChartVersion
}

//HTTPClient defines a behaviour for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
