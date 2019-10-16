package wsgatherer

import (
	"encoding/json"
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

func ParseJWT(tokenString string) []byte {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSecret(), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		fmt.Println("Stream Id from JWT: ", claims["stream_id"])
	} else {
		fmt.Println(err)
	}

	res, _ := json.Marshal(claims)
	fmt.Println("Data from JWT to JSON", res)

	return res
}
