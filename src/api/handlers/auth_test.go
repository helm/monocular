package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/gorilla/sessions"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

var dbSession = models.NewMockSession(models.MockDBConfig{})

func TestNewAuthHandlers(t *testing.T) {
	tests := []struct {
		name          string
		setSigningKey bool
		wantErr       bool
	}{
		{"no signing key", false, true},
		{"signing key", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setSigningKey {
				os.Setenv("MONOCULAR_AUTH_SIGNING_KEY", "secret")
				defer os.Unsetenv("MONOCULAR_AUTH_SIGNING_KEY")
			}
			got, err := NewAuthHandlers(dbSession)
			if tt.wantErr {
				assert.ExistsErr(t, err, "NewAuthHandlers error")
				assert.Nil(t, got, "returned AuthHandlers")
				return
			}
			assert.Equal(t, got.signingKey, "secret", "signing key")
		})
	}
}

func TestAuthHandlers_InitiateOAuth(t *testing.T) {
	s := sessions.NewCookieStore([]byte("secret"))
	ah := &AuthHandlers{"secret", s, dbSession}
	os.Setenv("MONOCULAR_AUTH_GITHUB_CLIENT_ID", "clientid")
	defer os.Unsetenv("MONOCULAR_AUTH_GITHUB_CLIENT_ID")
	os.Setenv("MONOCULAR_AUTH_GITHUB_CLIENT_SECRET", "clientsecret")
	defer os.Unsetenv("MONOCULAR_AUTH_GITHUB_CLIENT_SECRET")

	req, err := http.NewRequest("GET", "/auth", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	ah.InitiateOAuth(res, req)

	assert.Equal(t, res.Code, http.StatusFound, "response code")

	session, err := s.Get(req, "ka_sess")
	assert.NoErr(t, err)
	_, ok := session.Values["state"]
	assert.True(t, ok, "expected state to be set in session")

	location := res.Header().Get("Location")
	assert.True(t, location != "", "expected Location header to be set")
}

func TestAuthHandlers_InitiateOAuthForbidden(t *testing.T) {
	s := sessions.NewCookieStore([]byte("secret"))
	ah := &AuthHandlers{"secret", s, dbSession}
	req, err := http.NewRequest("GET", "/auth", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	ah.InitiateOAuth(res, req)
	assert.Equal(t, res.Code, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.NotNil(t, httpBody.Message, "error response")
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.True(t, strings.Contains(*httpBody.Message, "auth service not enabled"), "error message in HTTP body data")
}

func TestAuthHandlers_GithubCallbackStateMismatch(t *testing.T) {
	s := sessions.NewCookieStore([]byte("secret"))
	ah := &AuthHandlers{"secret", s, dbSession}
	os.Setenv("MONOCULAR_AUTH_GITHUB_CLIENT_ID", "clientid")
	defer os.Unsetenv("MONOCULAR_AUTH_GITHUB_CLIENT_ID")
	os.Setenv("MONOCULAR_AUTH_GITHUB_CLIENT_SECRET", "clientsecret")
	defer os.Unsetenv("MONOCULAR_AUTH_GITHUB_CLIENT_SECRET")

	req, err := http.NewRequest("GET", "/auth/github/callback", nil)
	assert.NoErr(t, err)
	state := "state"
	session, err := s.Get(req, "ka_sess")
	assert.NoErr(t, err)
	session.Values["state"] = state

	res := httptest.NewRecorder()
	ah.GithubCallback(res, req)

	assert.Equal(t, res.Code, http.StatusBadRequest, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.NotNil(t, httpBody.Message, "error response")
	assert.Equal(t, *httpBody.Code, int64(http.StatusBadRequest), "response code in HTTP body data")
	assert.True(t, strings.Contains(*httpBody.Message, "no state match"), "error message in HTTP body data")
}

func TestAuthHandlers_GithubCallbackForbidden(t *testing.T) {
	s := sessions.NewCookieStore([]byte("secret"))
	ah := &AuthHandlers{"secret", s, dbSession}

	req, err := http.NewRequest("GET", "/auth/github/callback", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	ah.GithubCallback(res, req)

	assert.Equal(t, res.Code, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.NotNil(t, httpBody.Message, "error response")
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.True(t, strings.Contains(*httpBody.Message, "auth service not enabled"), "error message in HTTP body data")
}

func TestAuthHandlers_Logout(t *testing.T) {
	s := sessions.NewCookieStore([]byte("secret"))
	ah := &AuthHandlers{"secret", s, dbSession}

	req, err := http.NewRequest("DELETE", "/auth/logout", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	ah.Logout(res, req)

	assert.Equal(t, res.Code, http.StatusOK, "response code")
	header := res.Header().Get("Set-Cookie")
	assert.Equal(t, header, "ka_auth=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT", "header")
}

func Test_tokenExpiration(t *testing.T) {
	expectedHour := (time.Now().Hour() + 2) % 24
	assert.Equal(t, tokenExpiration().Hour(), expectedHour, "hour")
}

func Test_randomStr(t *testing.T) {
	s1 := randomStr()
	s2 := randomStr()
	assert.True(t, s1 != s2, "expected %s to not equal %s", s1, s2)
}

func Test_errorResponse(t *testing.T) {
	message := "error message"
	res := httptest.NewRecorder()
	errorResponse(res, http.StatusBadRequest, message)
	assert.Equal(t, res.Code, http.StatusBadRequest, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusBadRequest), "response code in body")
	assert.Equal(t, *httpBody.Message, message, "error message in body")
}
