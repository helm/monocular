package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"gopkg.in/mgo.v2/bson"
)

// User describes a user
type User struct {
	ID      bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	Login   string          `json:"login"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	starred []bson.ObjectId `bson:"starred"`
}

// UserClaims describes a JWT claim for a user
type UserClaims struct {
	*User
	jwt.StandardClaims
}

// Context key type for request contexts
type contextKey int

// UserKey is the context key for the User data in the request context
const UserKey contextKey = 0

// UsersCollection is the name of the users collection
const UsersCollection = "users"

// CreateUser takes a User object and saves it to the database
// Users with the same login id are updated
func CreateUser(db datastore.Database, user *User) error {
	c := db.C(UsersCollection)
	_, err := c.Upsert(bson.M{"login": user.Login}, user)
	return err
}

// GetUserByLogin finds a user given the login id
func GetUserByLogin(db datastore.Database, login string) (*User, error) {
	c := db.C(UsersCollection)
	var user User
	err := c.Find(bson.M{"login": login}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
