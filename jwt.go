package golangmodule

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	golangmodule "github.com/rahmathidayat72/golang-module"
	"github.com/sirupsen/logrus"
)

type MetaToken struct {
	ID  int    `json:"id"`
	Exp string `json:"exp"`
}

type AccessToken struct {
	Claims MetaToken
}

func Sign(data map[string]interface{}) (string, time.Time, error) {
	// Menetapkan waktu kedaluwarsa token secara hardcode
	expiryTime := time.Now().UTC().Add(time.Hour * 48) // Waktu kedaluwarsa 48 jam di UTC

	claims := jwt.MapClaims{}
	claims["exp"] = expiryTime.Unix()

	for key, value := range data {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	// Format waktu kedaluwarsa sebagai string untuk kemudahan penggunaan
	//expiryTimeString := expiryTime.Format(time.RFC3339)

	return accessToken, expiryTime, nil
}

func VerifyTokenHeader(requestToken string) (MetaToken, error) {

	token, err := jwt.Parse((requestToken), func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return MetaToken{}, golangmodule.BuildResponse(http.StatusUnauthorized, "Token tidak valid")
	}

	claimToken := DecodeToken(token)
	return claimToken.Claims, nil
}

func VerifyToken(accessToken string) (*jwt.Token, error) {
	jwtSecretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return nil, golangmodule.BuildResponse(http.StatusUnauthorized, "Token tidak valid")
	}
		
	return token, nil
}

func DecodeToken(accessToken *jwt.Token) AccessToken {
	var token AccessToken
	stringify, err := json.Marshal(&accessToken)
	if err != nil {
		return token
	}
	err = json.Unmarshal(stringify, &token)
	if err != nil {
		return token
	}
	return token
}

// Fungsi untuk mengambil token dari header Authorization
func GetTokenFromAuthorizationHeader(authorizationHeader string) string {
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
