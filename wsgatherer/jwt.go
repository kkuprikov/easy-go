// Package wsgatherer - this file contains JWT logic
package wsgatherer

import (
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

const secretPath = ".secret"

func hmacSecret() []byte {
	res, err := ioutil.ReadFile(secretPath)
	if err != nil {
		fmt.Println("Could not read secret file: ", secretPath)
	}

	return res
}

func parseJWT(tokenString string) ([]byte, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret(), nil
	})

	if err != nil {
		fmt.Println("JWT parsing error: ", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid JWT token")
	}

	res, err := json.Marshal(claims)

	fmt.Println("Data from JWT to JSON", string(res))

	if err != nil {
		fmt.Println("JSON marshalling error: ", err)
		return nil, err
	}

	return res, nil
}
