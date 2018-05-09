package auth

import (
	"net/http"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	mySigningKey               = []byte("SuperSecret")
	DefaultKeyFunc jwt.Keyfunc = func(t *jwt.Token) ([]byte, error) { return mySigningKey, nil }
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: defaultKeyFunc,
	SigningMethod:       &jwt.SigningMethodHS256{},
})

func Validate(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ErrEmptyToken
	}
	token := strings.Replace(strings.TrimSpace(authHeader), "Bearer ", "", 1)
	_, err := jwt.Parse(token, defaultKeyFunc)
	if err != nil {
		return err
	}
	return nil
}

func Generate(user User) (string, error) {
	token := jwt.New(&jwt.SigningMethodHS256{})

	token.Claims["name"] = user.Name
	token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString(mySigningKey)
}
