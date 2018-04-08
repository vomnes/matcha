package lib

import (
	"errors"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
)

// AnalyseJWT check the JWT validity and get the data from the JWT
func AnalyseJWT(strToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return JWTSecret, nil
	})
	if err != nil {
		if ve, yes := err.(*jwt.ValidationError); yes {
			if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
				return nil, errors.New("Token expired")
			}
		}
		log.Println(PrettyError("[JWT] Not a valid JSON Web Token - " + err.Error()))
		return nil, errors.New("Not a valid JSON Web Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Not a valid token")
	}
	return claims, nil
}
