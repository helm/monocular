package middleware

import (
	"context"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/urfave/negroni"

	"github.com/kubernetes-helm/monocular/src/api/handlers/renderer"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

// AuthGate implements middleware to check if the user is logged in before continuing
func AuthGate() negroni.HandlerFunc {
	enabled := true
	signingKey, err := config.GetAuthSigningKey()
	if err != nil {
		enabled = false
	}
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		if !enabled {
			next(w, req)
			return
		}
		cookie, err := req.Cookie("ka_auth")
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})
		if err != nil {
			unauthorizedResponse(w)
			return
		}

		if claims, ok := token.Claims.(*models.UserClaims); ok && token.Valid {
			ctx := context.WithValue(req.Context(), models.UserKey, *claims)
			next(w, req.WithContext(ctx))
		} else {
			unauthorizedResponse(w)
		}
	}
}

func unauthorizedResponse(w http.ResponseWriter) {
	renderer.Render.JSON(w, http.StatusUnauthorized, swaggermodels.Error{
		Code:    pointerto.Int64(int64(http.StatusUnauthorized)),
		Message: pointerto.String("not logged in"),
	})
}
