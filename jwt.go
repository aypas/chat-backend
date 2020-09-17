package main
import (
	"os"
	"github.com/dgrijalva/jwt-go"
)

type CustomClaim struct {
	User_Id 	int 		`json:"user_id"`
	Exp 		float32 	`json:"exp"`
	Jti 		string 		`json:"jti"`
	Token_Type  string 		`json:"token_type"`
}

func loadSecret() []byte {
	s, b := os.LookupEnv("secret")
	if !b {
		panic("secret key not found")
	}
	return []byte(s)
}

func validate(tokenS string) bool {
	t, err := jwt.Parse(tokenS, func(_ *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return false
	}
	return t.Valid
}