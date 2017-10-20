package datastore

import (
	"errors"

	"gopkg.in/mgo.v2"
)

// Config configures the database connection
type Config struct {
	Host     string
	Database string
}

// Session is an interface for a MongoDB session
type Session interface {
	DB() (Database, func())
}

// Database is an interface for accessing a MongoDB database
type Database interface {
	C(name string) Collection
}

// Collection is an interface for accessing a MongoDB collection
type Collection interface {
	Find(query interface{}) Query
	// Count() (n int, err error)
	// Insert(docs ...interface{}) error
	Remove(selector interface{}) error
	UpdateId(selector, update interface{}) error
	Upsert(selector, update interface{}) (*mgo.ChangeInfo, error)
}

// Query is an interface for querying a MongoDB collection
type Query interface {
	All(result interface{}) error
	One(result interface{}) error
}

// mgoSession wraps an mgo.Session and implements Session
type mgoSession struct {
	conf Config
	*mgo.Session
}

func (s *mgoSession) DB() (Database, func()) {
	copy := s.Session.Copy()
	closer := func() { copy.Close() }
	return &mgoDatabase{copy.DB(s.conf.Database)}, closer
}

// mgoDatabase wraps an mgo.Database and implements Database
type mgoDatabase struct {
	*mgo.Database
}

func (d *mgoDatabase) C(name string) Collection {
	return &mgoCollection{d.Database.C(name)}
}

// mgoCollection wraps an mgo.Collection and implements Collection
type mgoCollection struct {
	*mgo.Collection
}

func (c *mgoCollection) Find(query interface{}) Query {
	return c.Collection.Find(query)
}

// NewSession initializes a MongoDB connection to the given host
func NewSession(conf Config) (Session, error) {
	session, err := mgo.Dial(conf.Host)
	if err != nil {
		return nil, errors.New("unable to connect to MongoDB")
	}
	session.SetMode(mgo.Monotonic, true)
	return &mgoSession{conf, session}, nil
}
