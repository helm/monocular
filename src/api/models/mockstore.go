package models

import (
	"errors"

	"github.com/kubernetes-helm/monocular/src/api/datastore"
	mgo "gopkg.in/mgo.v2"
)

// MockDBConfig describes configuration for the mock db session
type MockDBConfig struct {
	// If WantErr is set to true, all collection methods and query methods will return errors
	WantErr bool
	// If Empty is set to true, the Query methods will not affect the result parameter
	Empty bool
}

// mockSession acts as a mock datastore.Session
type mockSession struct {
	c MockDBConfig
}

func (s mockSession) DB() (datastore.Database, func()) {
	return mockDatabase{s.c}, func() {}
}

// mockDatabase acts as a mock datastore.Database
type mockDatabase struct {
	c MockDBConfig
}

func (d mockDatabase) C(name string) datastore.Collection {
	return mockCollection{d.c, name}
}

// mockCollection acts as a mock datastore.Collection
type mockCollection struct {
	c    MockDBConfig
	name string
}

func (c mockCollection) Find(query interface{}) datastore.Query {
	return mockQuery{c.c, c.name}
}

func (c mockCollection) Upsert(selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	if c.c.WantErr {
		return nil, errors.New("query error")
	}
	return nil, nil
}

func (c mockCollection) UpdateId(selector, update interface{}) error {
	if c.c.WantErr {
		return errors.New("query error")
	}
	return nil
}

func (c mockCollection) Remove(selector interface{}) error {
	if c.c.WantErr {
		return errors.New("query error")
	}
	return nil
}

// mockQuery acts as a mock datastore.Query
type mockQuery struct {
	c          MockDBConfig
	collection string
}

// Return values

// MockRepos contains the mocked repos returned by the database
var MockRepos = OfficialRepos

// MockUser contains the mock user returned by the database
var MockUser = &User{Name: "Rick Sanchez", Email: "rick@sanchez.com"}

func (q mockQuery) All(result interface{}) error {
	if q.c.Empty {
		return nil
	}
	if q.c.WantErr {
		return errors.New("query error")
	}
	switch q.collection {
	case ReposCollection:
		*result.(*[]*Repo) = MockRepos
	default:
		panic(q.collection + " mock not implemented")
	}
	return nil
}

func (q mockQuery) One(result interface{}) error {
	if q.c.Empty {
		return nil
	}
	if q.c.WantErr {
		return errors.New("query error")
	}
	switch q.collection {
	case ReposCollection:
		*result.(*Repo) = *MockRepos[0]
	case UsersCollection:
		*result.(*User) = *MockUser
	default:
		panic(q.collection + " mock not implemented")
	}
	return nil
}

// NewMockSession returns a mocked Session
func NewMockSession(config MockDBConfig) datastore.Session {
	return mockSession{config}
}
