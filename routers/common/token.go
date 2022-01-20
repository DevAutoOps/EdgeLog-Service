package common

import (
	"encoding/json"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const USERNAME = "admin"
const PASSWORD = "123456"
const randomCharRange = "abcdefghijklmnopqestuvwxyzABCDEFGHIJKLMNOPQESTUVWXYZ"

func init() {
	for i := 0; i < len(encryptionKey); i++ {
		encryptionKey[i] = GetRandomByte()
	}
}

var (
	// 16  byte  - AES-128
	// 24  byte  - AES-192
	// 32  byte  - AES-256
	encryptionKey = make([]byte, 32)
	TokenMap      sync.Map
)

func GetRandomString(len int) string {
	var result string
	for i := 0; i < 8; i++ {
		result += string(randomCharRange[rand.Intn(52)])
	}
	return result
}

func GetRandomByte() byte {
	return randomCharRange[rand.Intn(52)]
}

type Token struct {
	Expire int64
}

func (t Token) ToBytes() []byte {
	b, err := json.Marshal(t)
	if err != nil {
		return nil
	}
	return b
}

func toToken(b []byte) (Token, error) {
	result := Token{}
	err := json.Unmarshal(b, &result)
	return result, err
}

func stringToToken(s string) (token Token, err error) {
	aesEncryptResult := Base64Decode(s)
	aesDecryptResult, err := AesDecrypt([]byte(aesEncryptResult), encryptionKey)
	if err != nil {
		return
	}
	token, err = toToken(aesDecryptResult)
	return
}

func tokenToString(t Token) (string, error) {
	aesEncryptResult, err := AesEncrypt(t.ToBytes(), encryptionKey)
	if err != nil {
		return "", err
	}
	base64EncodeResult := Base64Encode(aesEncryptResult)
	return base64EncodeResult, nil
}

func GetOrAddToken(tokenKey string, token Token) (string, error) {
	load, ok := TokenMap.Load(tokenKey)
	if ok {
		return load.(string), nil
	}
	tokenStr, err := tokenToString(token)
	if err != nil {
		return "", err
	}
	TokenMap.Store(tokenKey, tokenStr)
	return tokenStr, nil
}

func TokenMid(c *gin.Context) {
	tokenHeader := c.Request.Header.Get("token")
	if tokenHeader == "" {
		c.JSON(401, gin.H{
			"message": "need 'token' Header",
		})
		c.Abort()
		return
	}
	load, ok := TokenMap.Load(GetTokenKey(c))
	if !ok {
		c.JSON(401, gin.H{
			"message": "unknown client",
		})
		c.Abort()
		return
	}
	tokenStr := load.(string)
	if tokenHeader != tokenStr {
		c.JSON(401, gin.H{
			"message": "wrongful token",
		})
		c.Abort()
		return
	}
	token, err := stringToToken(tokenStr)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		c.Abort()
	}
	if time.Now().After(time.Unix(token.Expire, 0)) {
		TokenMap.Delete(GetTokenKey(c))
		c.JSON(401, gin.H{
			"message": "token expire",
		})
		c.Abort()
	}
	c.Next()
}

func GetTokenKey(c *gin.Context) string {
	tokenKey := c.Request.RemoteAddr
	if strings.Contains(tokenKey, ":") {
		tokenKey = tokenKey[:strings.Index(tokenKey, ":")]
	}
	return tokenKey
}
