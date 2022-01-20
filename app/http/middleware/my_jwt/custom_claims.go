package my_jwt

import "github.com/dgrijalva/jwt-go"

//  custom jwt Declaration field information for + Standard field ， Reference address ：https://blog.csdn.net/codeSquare/article/details/99288718
type CustomClaims struct {
	UserId int64  `json:"user_id"`
	Name   string `json:"user_name"`
	Phone  string `json:"phone"`
	jwt.StandardClaims
}
