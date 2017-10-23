package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/arschles/assert"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

var handler = func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuthGate(t *testing.T) {
	tests := []struct {
		name   string
		cookie http.Cookie
		want   int
	}{
		{"no cookie", http.Cookie{}, http.StatusUnauthorized},
		{"unparseable JWT", http.Cookie{Name: "ka_auth", Value: "unparseable", Path: "/"}, http.StatusUnauthorized},
		{"expired JWT", generateJWT(generateClaims(time.Now().Add(-time.Hour * 2))), http.StatusUnauthorized},
		{"valid JWT", generateJWT(generateClaims(time.Now().Add(time.Hour * 2))), http.StatusOK},
	}

	os.Setenv("MONOCULAR_AUTH_SIGNING_KEY", "secret")
	defer os.Unsetenv("MONOCULAR_AUTH_SIGNING_KEY")
	for _, tt := range tests {
		req, err := http.NewRequest("GET", "/auth-gate", nil)
		assert.NoErr(t, err)
		req.AddCookie(&tt.cookie)
		res := httptest.NewRecorder()
		AuthGate()(res, req, handler)
		assert.Equal(t, res.Code, tt.want, tt.name)
	}
}

func TestAuthGateDisabled(t *testing.T) {
	req, err := http.NewRequest("GET", "/auth-gate", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	AuthGate()(res, req, handler)
	assert.Equal(t, res.Code, http.StatusOK, "response code")
}

func generateClaims(expiresAt time.Time) jwt.Claims {
	return models.UserClaims{
		User: &models.User{Name: "Jon Snow", Email: "jonsnow@winteriscoming.io"},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			Issuer:    "lyanna.stark",
		},
	}
}

func generateJWT(claims jwt.Claims) http.Cookie {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("secret"))
	return http.Cookie{Name: "ka_auth", Value: signedToken, Path: "/", HttpOnly: true}
}

func Test_unauthorizedResponse(t *testing.T) {
	res := httptest.NewRecorder()
	unauthorizedResponse(res)
	assert.Equal(t, res.Code, http.StatusUnauthorized, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusUnauthorized), "response code in body")
	assert.Equal(t, *httpBody.Message, "not logged in", "error message in body")
}
