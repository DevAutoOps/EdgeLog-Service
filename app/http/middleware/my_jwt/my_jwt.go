package my_jwt

import (
	"edgelog/app/global/my_errors"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

//  Create a using the factory  JWT  structural morphology 
func CreateMyJWT(signKey string) *JwtSign {
	if len(signKey) <= 0 {
		signKey = "edgelog"
	}
	return &JwtSign{
		[]byte(signKey),
	}
}

//  Define a  JWT Signature verification   structural morphology 
type JwtSign struct {
	SigningKey []byte
}

// CreateToken  Generate a token
func (j *JwtSign) CreateToken(claims CustomClaims) (string, error) {
	//  generate jwt Formatted header、claims  part 
	tokenPartA := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//  Continue adding secret key values ， Generate the last part 
	return tokenPartA.SignedString(j.SigningKey)
}

//  analysis Token
func (j *JwtSign) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if token == nil {
		return nil, errors.New(my_errors.ErrorsTokenInvalid)
	}
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New(my_errors.ErrorsTokenMalFormed)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New(my_errors.ErrorsTokenNotActiveYet)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				//  If  TokenExpired , Just expired （ The format is correct ）， We think he is effective ， Next, you can allow the refresh operation 
				token.Valid = true
				goto labelHere
			} else {
				return nil, errors.New(my_errors.ErrorsTokenInvalid)
			}
		}
	}
labelHere:
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New(my_errors.ErrorsTokenInvalid)
	}
}

//  to update token
func (j *JwtSign) RefreshToken(tokenString string, extraAddSeconds int64) (string, error) {

	if CustomClaims, err := j.ParseToken(tokenString); err == nil {
		CustomClaims.ExpiresAt = time.Now().Unix() + extraAddSeconds
		return j.CreateToken(*CustomClaims)
	} else {
		return "", err
	}
}
