package tools

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"runtime"
	"strconv"
)

func CompareHashAndPassword(e string, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Assert  Conditional assertion
//  When the assertion condition is   false   Time trigger  panic
//  The following code will not be executed for the current request ， And return the error information and error code in the specified format
func Assert(condition bool, msg string, code ...int) {
	if !condition {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}

// HasError  False assertion
//  When  error  Not for  nil  Time trigger  panic
//  The following code will not be executed for the current request ， And return the error information and error code in the specified format
//  if  msg  Empty ， The default is  error  Content in
func HasError(err error, msg string, code ...int) {
	if err != nil {
		statusCode := 200
		if len(code) > 0 {
			statusCode = code[0]
		}
		if msg == "" {
			msg = err.Error()
		}
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%v error: %#v", file, line, err)
		panic("CustomError#" + strconv.Itoa(statusCode) + "#" + msg)
	}
}

// GenerateMsgIDFromContext  generate msgID
func GenerateMsgIDFromContext(c *gin.Context) string {
	var msgID string
	data, ok := c.Get("msgID")
	if !ok {
		msgID = uuid.New().String()
		c.Set("msgID", msgID)
		return msgID
	}
	msgID = cast.ToString(data)
	return msgID
}

// GetOrm  obtain orm connect
func GetOrm(c *gin.Context) (*gorm.DB, error) {
	msgID := GenerateMsgIDFromContext(c)
	idb, exist := c.Get("db")
	if !exist {
		return nil, errors.New(fmt.Sprintf("msgID[%s], db connect not exist", msgID))
	}
	switch idb.(type) {
	case *gorm.DB:
		// Add operation
		return idb.(*gorm.DB), nil
	default:
		return nil, errors.New(fmt.Sprintf("msgID[%s], db connect not exist", msgID))
	}
}
