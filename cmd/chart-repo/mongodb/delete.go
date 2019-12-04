/*
Copyright (c) 2018 The Helm Authors

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

package mongodb

import (
	"os"

	"github.com/kubeapps/common/datastore"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//Delete deletes a repository from a mongoDB datastore
func Delete(cmd *cobra.Command, args []string) {
	mongoURL, err := cmd.Flags().GetString("mongo-url")
	if err != nil {
		logrus.Fatal(err)
	}
	mongoDB, err := cmd.Flags().GetString("mongo-database")
	if err != nil {
		logrus.Fatal(err)
	}
	mongoUser, err := cmd.Flags().GetString("mongo-user")
	if err != nil {
		logrus.Fatal(err)
	}
	mongoPW := os.Getenv("MONGO_PASSWORD")
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		logrus.Fatal(err)
	}
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	mongoConfig := datastore.Config{URL: mongoURL, Database: mongoDB, Username: mongoUser, Password: mongoPW}
	dbSession, err := datastore.NewSession(mongoConfig)
	if err != nil {
		logrus.Fatalf("Can't connect to mongoDB: %v", err)
	}
	if err = deleteRepo(dbSession, args[0]); err != nil {
		logrus.Fatalf("Can't delete chart repository %s from database: %v", args[0], err)
	}

	logrus.Infof("Successfully deleted the chart repository %s from database", args[0])
}
