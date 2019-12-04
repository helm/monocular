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

package foundationdb

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/helm/monocular/cmd/chart-repo/common"
	"github.com/helm/monocular/cmd/chart-repo/utils"

	"github.com/disintegration/imaging"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	chartCollection       = "charts"
	repositoryCollection  = "repos"
	chartFilesCollection  = "files"
	defaultTimeoutSeconds = 10
	additionalCAFile      = "/usr/local/share/ca-certificates/ca.crt"
)

var netClient common.HTTPClient = &http.Client{}

func init() {
	var err error
	netClient, err = common.InitNetClient(additionalCAFile, defaultTimeoutSeconds)
	if err != nil {
		log.Fatal(err)
	}
}

// SyncRepo Syncing is performed in the following steps:
// 1. Update database to match chart metadata from index
// 2. Concurrently process icons for charts (concurrently)
// 3. Concurrently process the README and values.yaml for the latest chart version of each chart
// 4. Concurrently process READMEs and values.yaml for historic chart versions
//
// These steps are processed in this way to ensure relevant chart data is
// imported into the database as fast as possible. E.g. we want all icons for
// charts before fetching readmes for each chart and version pair.
func syncRepo(dbClient Client, dbName, repoName, repoURL string, authorizationHeader string) error {

	db, closer := dbClient.Database(dbName)
	defer closer()

	url, err := common.ParseRepoURL(repoURL)
	if err != nil {
		log.WithFields(log.Fields{"url": repoURL}).WithError(err).Error("failed to parse URL")
		return err
	}

	log.Debugf("Checking database connection and readiness...")
	collection := db.Collection("numbers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159}, options.InsertOne())
	if err != nil {
		log.Fatalf("Database readiness test failed: %v", err)
		cancel()
		return err
	}
	id := res.InsertedID
	cancel()
	log.Debugf("Database connection test successful.")
	log.Debugf("Inserted a test document to test collection with ID: %v", id)

	r := common.Repo{Name: repoName, URL: url.String(), AuthorizationHeader: authorizationHeader}
	repoBytes, err := common.FetchRepoIndex(r, netClient)
	if err != nil {
		return err
	}

	repoChecksum, err := common.GetSha256(repoBytes)
	if err != nil {
		return err
	}

	// Check if the repo has been already processed
	if repoAlreadyProcessed(db, repoName, repoChecksum) {
		log.WithFields(log.Fields{"url": repoURL}).Info("Skipping repository since there are no updates")
		return nil
	}

	index, err := common.ParseRepoIndex(repoBytes)
	if err != nil {
		return err
	}

	charts := common.ChartsFromIndex(index, r)
	log.Debugf("%v Charts in index of repo: %v", len(charts), repoURL)
	if len(charts) == 0 {
		return errors.New("no charts in repository index")
	}

	err = importCharts(db, dbName, charts)
	if err != nil {
		return err
	}

	// Process 10 charts at a time
	numWorkers := 10
	iconJobs := make(chan common.Chart, numWorkers)
	chartFilesJobs := make(chan common.ImportChartFilesJob, numWorkers)
	var wg sync.WaitGroup

	log.Debugf("starting %d workers", numWorkers)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go importWorker(db, &wg, iconJobs, chartFilesJobs)
	}

	// Enqueue jobs to process chart icons
	for _, c := range charts {
		iconJobs <- c
	}
	// Close the iconJobs channel to signal the worker pools to move on to the
	// chart files jobs
	close(iconJobs)

	// Iterate through the list of charts and enqueue the latest chart version to
	// be processed. Append the rest of the chart versions to a list to be
	// enqueued later
	var toEnqueue []common.ImportChartFilesJob
	for _, c := range charts {
		chartFilesJobs <- common.ImportChartFilesJob{Name: c.Name, Repo: c.Repo, ChartVersion: c.ChartVersions[0]}
		for _, cv := range c.ChartVersions[1:] {
			toEnqueue = append(toEnqueue, common.ImportChartFilesJob{Name: c.Name, Repo: c.Repo, ChartVersion: cv})
		}
	}

	// Enqueue all the remaining chart versions
	for _, cfj := range toEnqueue {
		chartFilesJobs <- cfj
	}
	// Close the chartFilesJobs channel to signal the worker pools that there are
	// no more jobs to process
	close(chartFilesJobs)

	// Wait for the worker pools to finish processing
	wg.Wait()

	// Update cache in the database
	if err = updateLastCheck(db, repoName, repoChecksum, time.Now()); err != nil {
		return err
	}
	log.WithFields(log.Fields{"url": repoURL}).Info("Stored repository update in cache")

	return nil
}

func repoAlreadyProcessed(db Database, repoName string, checksum string) bool {
	lastCheck := &common.RepoCheck{}
	filter := bson.M{"_id": repoName}
	err := db.Collection(repositoryCollection).FindOne(context.Background(), filter, lastCheck, options.FindOne())
	return err == nil && checksum == lastCheck.Checksum
}

func updateLastCheck(db Database, repoName string, checksum string, now time.Time) error {
	selector := bson.M{"_id": repoName}
	update := bson.M{"$set": bson.M{"last_update": now, "checksum": checksum}}
	_, err := db.Collection(repositoryCollection).UpdateOne(context.Background(), selector, update, options.Update())
	return err
}

func deleteRepo(dbClient Client, dbName, repoName string) error {
	db, closer := dbClient.Database(dbName)
	defer closer()
	collection := db.Collection(chartCollection)
	filter := bson.M{
		"repo.name": repoName,
	}
	deleteResult, err := collection.DeleteMany(context.Background(), filter, options.Delete())
	if err != nil {
		log.Debugf("Error occurred during delete repo (deleting charts from index). Err: %v, Result: %v", err, deleteResult)
		return err
	}
	log.Debugf("Repo delete (delete charts from index) result: %v charts deleted", deleteResult.DeletedCount)

	collection = db.Collection(chartFilesCollection)
	deleteResult, err = collection.DeleteMany(context.Background(), filter, options.Delete())
	if err != nil {
		log.Debugf("Error occurred during delete repo (deleting chart files from index). Err: %v, Result: %v", err, deleteResult)
		return err
	}
	log.Debugf("Repo delete (delete chart files from index) result: %v chart files deleted.", deleteResult.DeletedCount)
	collection = db.Collection(repositoryCollection)
	deleteResult, err = collection.DeleteMany(context.Background(), filter, options.Delete())
	if err != nil {
		log.Debugf("Error occurred during delete repo (deleting repositories from index). Err: %v, Result: %v", err, deleteResult)
		return err
	}
	log.Debugf("Repo delete (delete chart files from index) result: %v repositories deleted.", deleteResult.DeletedCount)

	return err
}

func importCharts(db Database, dbName string, charts []common.Chart) error {
	var operations []mongo.WriteModel
	var chartIDs []string
	for _, c := range charts {
		operation := mongo.NewUpdateOneModel()
		chartIDs = append(chartIDs, c.ID)
		// charts to upsert - pair of filter, chart
		operation.SetFilter(bson.M{
			"_id": c.ID,
		})

		chartBSON, err := bson.Marshal(&c)
		var doc bson.M
		bson.Unmarshal(chartBSON, &doc)
		delete(doc, "_id")

		if err != nil {
			log.Debugf("Error marshalling chart to BSON: %v. Skipping this chart.", err)
		} else {
			update := doc
			operation.SetUpdate(update)
			operation.SetUpsert(true)
			operations = append(operations, operation)
		}
		log.Debugf("Adding chart insert operation for chart: %v", c.ID)
	}

	//Must use bulk write for array of filters
	collection := db.Collection(chartCollection)
	updateResult, err := collection.BulkWrite(
		context.Background(),
		operations,
		options.BulkWrite(),
	)

	//Set upsert flag and upsert the pairs here
	//Updates our index for charts that we already have and inserts charts that are new
	if err != nil {
		log.Debugf("Error occurred during chart import (upsert many). Err: %v", err)
		return err
	}
	log.Debugf("Upsert chart index success. %v documents inserted, %v documents upserted, %v documents modified", updateResult.InsertedCount, updateResult.UpsertedCount, updateResult.ModifiedCount)

	//Remove from our index, any charts that no longer exist
	filter := bson.M{
		"_id": bson.M{
			"$nin": chartIDs,
		},
		"repo.name": charts[0].Repo.Name,
	}
	deleteResult, err := collection.DeleteMany(context.Background(), filter, options.Delete())
	if err != nil {
		log.Debugf("Error occurred during chart import (delete many). Err: %v", err)
		return err
	}
	log.Debugf("Delete stale charts from index success. %v documents deleted.", deleteResult.DeletedCount)

	return err
}

func importWorker(db Database, wg *sync.WaitGroup, icons <-chan common.Chart, chartFiles <-chan common.ImportChartFilesJob) {
	defer wg.Done()
	for c := range icons {
		log.WithFields(log.Fields{"name": c.Name}).Debug("importing icon")
		if err := fetchAndImportIcon(db, c); err != nil {
			log.WithFields(log.Fields{"name": c.Name}).WithError(err).Error("failed to import icon")
		}
	}
	for j := range chartFiles {
		log.WithFields(log.Fields{"name": j.Name, "version": j.ChartVersion.Version}).Debug("importing readme and values")
		if err := fetchAndImportFiles(db, j.Name, j.Repo, j.ChartVersion); err != nil {
			log.WithFields(log.Fields{"name": j.Name, "version": j.ChartVersion.Version}).WithError(err).Error("failed to import files")
		}
	}
}

func fetchAndImportIcon(db Database, c common.Chart) error {
	if c.Icon == "" {
		log.WithFields(log.Fields{"name": c.Name}).Info("icon not found")
		return nil
	}

	req, err := http.NewRequest("GET", c.Icon, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", utils.UserAgent())
	if len(c.Repo.AuthorizationHeader) > 0 {
		req.Header.Set("Authorization", c.Repo.AuthorizationHeader)
	}

	res, err := netClient.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s", res.StatusCode, c.Icon)
	}

	b := []byte{}
	contentType := ""
	if strings.Contains(res.Header.Get("Content-Type"), "image/svg") {
		// if the icon is a SVG file simply read it
		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		contentType = res.Header.Get("Content-Type")
	} else {
		// if the icon is in any other format try to convert it to PNG
		orig, err := imaging.Decode(res.Body)
		if err != nil {
			log.WithFields(log.Fields{"name": c.Name}).WithError(err).Error("failed to decode icon")
			return err
		}

		// TODO: make this configurable?
		icon := imaging.Fit(orig, 160, 160, imaging.Lanczos)

		var buf bytes.Buffer
		imaging.Encode(&buf, icon, imaging.PNG)
		b = buf.Bytes()
		contentType = "image/png"
	}

	collection := db.Collection(chartCollection)
	//Update single icon
	update := bson.M{"$set": bson.M{"raw_icon": b, "icon_content_type": contentType}}
	filter := bson.M{"_id": c.ID}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update, options.Update())
	if err != nil {
		log.Debugf("Error occurred during chart icon import (update one). Err: %v, Result: %v", err, updateResult)
		return err
	}
	return err
}

func fetchAndImportFiles(db Database, name string, r common.Repo, cv common.ChartVersion) error {

	chartFilesID := fmt.Sprintf("%s/%s-%s", r.Name, name, cv.Version)
	//Check if we already have indexed files for this chart version and digest
	collection := db.Collection(chartFilesCollection)
	filter := bson.M{"_id": chartFilesID, "digest": cv.Digest}
	findResult := collection.FindOne(context.Background(), filter, &common.ChartFiles{}, options.FindOne())
	if findResult != mongo.ErrNoDocuments {
		log.WithFields(log.Fields{"name": name, "version": cv.Version}).Debug("skipping existing files")
		return nil
	}
	log.WithFields(log.Fields{"name": name, "version": cv.Version}).Debug("fetching files")

	url := common.ChartTarballURL(r, cv)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", utils.UserAgent())
	if len(r.AuthorizationHeader) > 0 {
		req.Header.Set("Authorization", r.AuthorizationHeader)
	}

	res, err := netClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// We read the whole chart into memory, this should be okay since the chart
	// tarball needs to be small enough to fit into a GRPC call (Tiller
	// requirement)
	gzf, err := gzip.NewReader(res.Body)
	if err != nil {
		return err
	}
	defer gzf.Close()

	tarf := tar.NewReader(gzf)

	readmeFileName := name + "/README.md"
	valuesFileName := name + "/values.yaml"
	schemaFileName := name + "/values.schema.json"
	filenames := []string{valuesFileName, readmeFileName, schemaFileName}

	files, err := common.ExtractFilesFromTarball(filenames, tarf)
	if err != nil {
		return err
	}

	chartFiles := common.ChartFiles{ID: chartFilesID, Repo: r, Digest: cv.Digest}
	if v, ok := files[readmeFileName]; ok {
		chartFiles.Readme = v
	} else {
		log.WithFields(log.Fields{"name": name, "version": cv.Version}).Info("README.md not found")
	}
	if v, ok := files[valuesFileName]; ok {
		chartFiles.Values = v
	} else {
		log.WithFields(log.Fields{"name": name, "version": cv.Version}).Info("values.yaml not found")
	}
	if v, ok := files[schemaFileName]; ok {
		chartFiles.Schema = v
	} else {
		log.WithFields(log.Fields{"name": name, "version": cv.Version}).Info("values.schema.json not found")
	}

	// inserts the chart files if not already indexed, or updates the existing
	// entry if digest has changed
	log.Debugf("Inserting chart files %v to collection: %v....", chartFilesID, chartFilesCollection)
	collection = db.Collection(chartFilesCollection)
	filter = bson.M{"_id": chartFilesID}
	chartBSON, err := bson.Marshal(&chartFiles)
	var doc bson.M
	bson.Unmarshal(chartBSON, &doc)
	delete(doc, "_id")
	update := bson.M{"$set": doc}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Debugf("Error occurred during chart files import (update one). Chart files : %v doc: %v  Err: %v", chartFiles, doc, err)
		return err
	}
	log.Debugf("Chart files import success. (update one) Upserted: %v Updated: %v", updateResult.UpsertedCount, updateResult.ModifiedCount)
	return nil
}

func database(client *mongo.Client, dbName string) (*mongo.Database, func()) {

	db := client.Database(dbName)
	return db, func() {
		err := client.Disconnect(context.Background())

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}
}
