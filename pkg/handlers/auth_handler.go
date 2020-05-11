package handlers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

// MyJWTSigningKey - just for test task, need to store in some config file for example
var MyJWTSigningKey = []byte("secret")

// just for test task, usually the login/encrypted_password should be in some storage
var validCredentials = map[string]string{
	"TestRssApp": "password",
}

type authHandler struct{}

// NewAuthHandler is a constructor for authHandler
func NewAuthHandler() *authHandler {
	return &authHandler{}
}

// Authenticate is a function which build new JWT token
// Route for it is POST /authenticate
// Incoming params are login and password which are strings
// The response is json of login and jwt token or some error message
func (h *authHandler) Authenticate(ctx *fasthttp.RequestCtx) {
	type successResponse struct {
		Login string `json:"login"`
		Token string `json:"token"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	login := string(ctx.PostArgs().Peek("login"))
	password := string(ctx.PostArgs().Peek("password"))

	if isValidCredentials(login, password) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"login": login,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, _ := token.SignedString(MyJWTSigningKey)

		respondWithJSON(ctx, 200, successResponse{
			Login: login,
			Token: tokenString,
		})
	} else {
		respondWithJSON(ctx, 401, errorResponse{
			Error: "NotAuthenticated - your credentials are incorrect",
		})
	}
}

func isValidCredentials(login string, password string) bool {
	for validLogin, validPassword := range validCredentials {
		if login == validLogin {
			if password == validPassword {
				return true
			}
		}
	}
	return false
}
