package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"

	"github.com/gorilla/sessions"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/handlers/renderer"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"golang.org/x/oauth2"
	oauth2Github "golang.org/x/oauth2/github"
)

type userClaims struct {
	Name  string
	Email string
	jwt.StandardClaims
}

// AuthHandlers defines handlers that provide authentication
type AuthHandlers struct {
	store sessions.Store
	conf  config.Configuration
}

// NewAuthHandlers takes a sessions.Store implementation and returns an AuthHandlers struct
func NewAuthHandlers(s sessions.Store, c config.Configuration) *AuthHandlers {
	return &AuthHandlers{s, c}
}

// InitiateOAuth initiatates an OAuth request
func (a *AuthHandlers) InitiateOAuth(w http.ResponseWriter, r *http.Request) {
	state := randomStr()
	session, _ := a.store.Get(r, "sess")
	session.Values["state"] = state
	session.Save(r, w)

	url := oauthConfig(a.conf, r.Host).AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// GithubCallback processes the OAuth callback from GitHub
func (a *AuthHandlers) GithubCallback(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "sess")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid session")
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		errorResponse(w, http.StatusBadRequest, "no state match - possible CSRF or cookies not enabled")
		return
	}

	tkn, err := oauthConfig(a.conf, r.Host).Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "unable to get access token")
		return
	}

	if !tkn.Valid() {
		errorResponse(w, http.StatusInternalServerError, "invalid access token retrieved")
		return
	}

	client := github.NewClient(oauthConfig(a.conf, r.Host).Client(oauth2.NoContext, tkn))

	user, _, err := client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "unable to retrieve user")
		return
	}

	claims := userClaims{
		*user.Name,
		*user.Email,
		jwt.StandardClaims{
			ExpiresAt: tokenExpiration().Unix(),
			Issuer:    r.Host,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(a.conf.SigningKey))
	jwtCookie := http.Cookie{Name: "ka_auth", Value: signedToken, Path: "/", Expires: tokenExpiration(), HttpOnly: true}

	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "error marshalling claims")
		return
	}
	claimsCookie := http.Cookie{Name: "ka_claims", Value: string(jsonClaims), Path: "/"}

	http.SetCookie(w, &jwtCookie)
	http.SetCookie(w, &claimsCookie)

	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func tokenExpiration() time.Time {
	return time.Now().Add(time.Hour * 2)
}

func oauthConfig(c config.Configuration, host string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.OAuthConfig.ClientID,
		ClientSecret: c.OAuthConfig.ClientSecret,
		Endpoint:     oauth2Github.Endpoint,
		RedirectURL:  "http://" + host + "/auth/github/callback",
		Scopes:       []string{"repo"},
	}
}

func randomStr() string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	renderer.Render.JSON(w, code, models.Error{Code: pointerto.Int64(int64(code)), Message: &message})
}
