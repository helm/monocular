package models

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"gopkg.in/mgo.v2/bson"
)

// User describes a user
type User struct {
	ID      bson.ObjectId   `json:"id" bson:"_id,omitempty"`
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
// Users with the same email are updated
func CreateUser(db datastore.Database, user *User) error {
	c := db.C(UsersCollection)
	_, err := c.Upsert(bson.M{"email": user.Email}, user)
	return err
}

// GetUserByEmail finds a user given an email
func GetUserByEmail(db datastore.Database, email string) (*User, error) {
	c := db.C(UsersCollection)
	var user User
	err := c.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// StarChart marks a chart as starred by the user
// It updates the `starred` property on the User to point to the Chart
// and the `stargazers` property on the Chart to point to the User
func (u *User) StarChart(db datastore.Database, chart *Chart) error {
	if err := db.C(UsersCollection).UpdateId(u.ID, bson.M{"$push": bson.M{"starred": chart.ID}}); err != nil {
		return err
	}
	if err := db.C(ChartsCollection).UpdateId(chart.ID, bson.M{"$push": bson.M{"stargazers": u.ID}}); err != nil {
		return err
	}
	return nil
}

// UnstarChart unmarks a chart as starred by the user
// It removes the reference from the `starred` property on the User pointing to the Chart
// and the `stargazers` property on the Chart pointing to the User
func (u *User) UnstarChart(db datastore.Database, chart *Chart) error {
	if err := db.C(UsersCollection).UpdateId(u.ID, bson.M{"$pull": bson.M{"starred": chart.ID}}); err != nil {
		return err
	}
	if err := db.C(ChartsCollection).UpdateId(chart.ID, bson.M{"$pull": bson.M{"stargazers": u.ID}}); err != nil {
		return err
	}
	return nil
}
