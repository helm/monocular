package datastore

import (
	"errors"
	"reflect"

	mgo "gopkg.in/mgo.v2"
)

// mockSession acts as a mock Session
type mockSession struct {
	mockResult interface{}
	wantErr    bool
}

func (s mockSession) DB() (Database, func()) {
	return mockDatabase{s.mockResult, s.wantErr}, func() {}
}

// mockDatabase acts as a mock Database
type mockDatabase struct {
	mockResult interface{}
	wantErr    bool
}

func (d mockDatabase) C(name string) Collection {
	return mockCollection{d.mockResult, d.wantErr}
}

// mockCollection acts as a mock Collection
type mockCollection struct {
	mockResult interface{}
	wantErr    bool
}

func (c mockCollection) Find(query interface{}) Query {
	return mockQuery{c.mockResult, c.wantErr}
}

func (c mockCollection) Upsert(selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	var err error
	if c.wantErr {
		err = errors.New("query error")
	}
	return nil, err
}

func (c mockCollection) UpdateId(selector, update interface{}) error {
	var err error
	if c.wantErr {
		err = errors.New("query error")
	}
	return err
}

func (c mockCollection) Remove(selector interface{}) error {
	if c.wantErr {
		return mgo.ErrNotFound
	}
	return nil
}

// mockQuery acts as a mock Query
type mockQuery struct {
	mockResult interface{}
	wantErr    bool
}

func (q mockQuery) All(result interface{}) error {
	if q.wantErr {
		return errors.New("query error")
	}
	copyTo(q.mockResult, result)
	return nil
}

func (q mockQuery) One(result interface{}) error {
	if q.wantErr {
		return errors.New("inexistant document")
	}
	copyTo(q.mockResult, result)
	return nil
}

func copyTo(src interface{}, dst interface{}) {
	// If src is nil, don't try to set
	if src == nil {
		return
	}
	srcV := reflect.ValueOf(src).Elem()
	dstV := reflect.ValueOf(dst).Elem()
	dstV.Set(srcV)
}

// NewMockSession returns a mocked Session
func NewMockSession(mockResult interface{}, wantErr bool) Session {
	return mockSession{mockResult, wantErr}
}
