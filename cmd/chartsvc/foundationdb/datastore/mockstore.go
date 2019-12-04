/*
Copyright (c) 2019

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

package datastore

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockDatabase acts as a mock datastore.Database
type mockDatabase struct {
	*mock.Mock
}

type mockClient struct {
	*mock.Mock
}

//NewMockClient returns a mocked FDB Document-Layer client
func NewMockClient(m *mock.Mock) Client {
	return mockClient{m}
}

// DB returns a mocked datastore.Database and empty closer function
func (c mockClient) Database(dbName string) (Database, func()) {

	db := mockDatabase{c.Mock}

	return db, func() {
	}
}

func (d mockDatabase) Collection(name string) Collection {
	return mockCollection{d.Mock}
}

// mockCollection acts as a mock datastore.Collection
type mockCollection struct {
	*mock.Mock
}

func (c mockCollection) BulkWrite(ctxt context.Context, operations []mongo.WriteModel, options *options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	args := c.Called(ctxt, operations, options)
	return args.Get(0).(*mongo.BulkWriteResult), args.Error(1)
}

func (c mockCollection) DeleteMany(ctxt context.Context, filter interface{}, options *options.DeleteOptions) (*mongo.DeleteResult, error) {
	args := c.Called(ctxt, filter, options)
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

func (c mockCollection) FindOne(ctxt context.Context, filter interface{}, result interface{}, options *options.FindOneOptions) error {
	args := c.Called(ctxt, filter, result, options)
	return args.Error(0)
}

func (c mockCollection) InsertOne(ctxt context.Context, document interface{}, options *options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	args := c.Called(ctxt, document, options)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (c mockCollection) UpdateOne(ctxt context.Context, filter interface{}, document interface{}, options *options.UpdateOptions) (*mongo.UpdateResult, error) {
	args := c.Called(ctxt, filter, document, options)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (c mockCollection) Find(ctxt context.Context, filter interface{}, result interface{}, options *options.FindOptions) error {
	args := c.Called(ctxt, filter, result, options)
	return args.Error(0)
}
