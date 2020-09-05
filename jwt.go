package main
import ("os"
		"fmt"
		"github.com/dgrijalva/jwt-go")

type CustomClaim struct {
	Exp float32 	`json:"exp"`
	Jti string 	`json:"jti"`
	Token_Type string `json:"token_type`
	User_Id int 	`json:"user_id"`
}

func loadSecret() []byte {
	s, b := os.LookupEnv("secret")
	if !b {
		panic("secret key not found")
	}
	return []byte(s)
}

func validate(tokenS string) bool {
	t, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		fmt.Println((*token).Method.Alg())
		return secret, nil
	})
	if err != nil {
		return false
	}
	return t.Valid
}