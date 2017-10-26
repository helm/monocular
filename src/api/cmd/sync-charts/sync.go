package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kubeapps/common/datastore"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var dbSession datastore.Session

var netClient httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [REPO NAME] [REPO URL]\n", os.Args[0])
		flag.PrintDefaults()
	}
	dbHost := flag.String("mongo-host", "localhost:27017", "MongoDB host")
	dbName := flag.String("mongo-database", "charts", "MongoDB database")
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	mongoConfig := datastore.Config{Host: *dbHost, Database: *dbName}
	var err error
	dbSession, err = datastore.NewSession(mongoConfig)
	if err != nil {
		log.WithFields(log.Fields{"host": *dbHost}).Fatal(err)
	}

	if err := sync(flag.Arg(0), flag.Arg(1)); err != nil {
		os.Exit(1)
	}
}

func sync(repoName, repoURL string) error {
	if _, err := url.ParseRequestURI(repoURL); err != nil {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("failed to parse URL")
		return err
	}

	req, err := http.NewRequest("GET", repoURL+"/index.yaml", nil)
	if err != nil {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("could not build repo index request")
		return err
	}
	res, err := netClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("error requesting repo index")
		return err
	}

	if res.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("error requesting repo index, are you sure this is a chart repository?")
		return errors.New("repo index request failed")
	}
	return nil
}
