package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hay-i/chronologger/models"
	"github.com/labstack/echo/v4"
)

func signToken(user models.User) (time.Time, string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   user.Username,
		ExpiresAt: expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))

	return expirationTime, signedToken, err
}

func setCookie(signedToken string, expiry time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = signedToken
	cookie.Expires = expiry
	cookie.HttpOnly = true // Make the cookie inaccessible to JavaScript running in the browser
	// https://github.com/hay-i/chronologger/issues/35
	// cookie.Secure = true
	c.SetCookie(cookie)
}
