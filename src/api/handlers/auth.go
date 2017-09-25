package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
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
)

type userClaims struct {
	Name  string
	Email string
	jwt.StandardClaims
}

// AuthHandlers defines handlers that provide authentication
type AuthHandlers struct {
	signingKey string
	store      sessions.Store
}

// NewAuthHandlers takes a sessions.Store implementation and returns an AuthHandlers struct
func NewAuthHandlers() (*AuthHandlers, error) {
	signingKey, err := config.GetAuthSigningKey()
	if err != nil {
		return nil, errors.New("no signing key, ensure MONOCULAR_AUTH_SIGNING_KEY is set")
	}
	s := sessions.NewCookieStore([]byte(signingKey))
	return &AuthHandlers{signingKey, s}, nil
}

// InitiateOAuth initiatates an OAuth request
func (a *AuthHandlers) InitiateOAuth(w http.ResponseWriter, r *http.Request) {
	oauthConfig, err := config.GetOAuthConfig(r.Host)
	if err != nil {
		errorResponse(w, http.StatusForbidden, "auth service not enabled: "+err.Error())
		return
	}
	state := randomStr()
	session, _ := a.store.Get(r, "ka_sess")
	session.Values["state"] = state
	session.Save(r, w)

	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// GithubCallback processes the OAuth callback from GitHub
func (a *AuthHandlers) GithubCallback(w http.ResponseWriter, r *http.Request) {
	oauthConfig, err := config.GetOAuthConfig(r.Host)
	if err != nil {
		errorResponse(w, http.StatusForbidden, "auth service not enabled: "+err.Error())
		return
	}

	session, err := a.store.Get(r, "ka_sess")
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid session")
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		errorResponse(w, http.StatusBadRequest, "no state match - possible CSRF or cookies not enabled")
		return
	}

	tkn, err := oauthConfig.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "unable to get access token")
		return
	}

	if !tkn.Valid() {
		errorResponse(w, http.StatusInternalServerError, "invalid access token retrieved")
		return
	}

	client := github.NewClient(oauthConfig.Client(oauth2.NoContext, tkn))

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
	signedToken, _ := token.SignedString([]byte(a.signingKey))
	jwtCookie := http.Cookie{Name: "ka_auth", Value: signedToken, Path: "/", Expires: tokenExpiration(), HttpOnly: true}

	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "error marshalling claims")
		return
	}
	claimsCookie := http.Cookie{Name: "ka_claims", Value: base64.StdEncoding.EncodeToString(jsonClaims), Path: "/"}

	http.SetCookie(w, &jwtCookie)
	http.SetCookie(w, &claimsCookie)

	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

// Logout clears the JWT token cookie
func (a *AuthHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{Name: "ka_auth", Value: "", Path: "/", Expires: time.Unix(1, 0)}
	http.SetCookie(w, &cookie)
}

func tokenExpiration() time.Time {
	return time.Now().Add(time.Hour * 2)
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
