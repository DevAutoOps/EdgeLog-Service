package interf

import "github.com/gin-gonic/gin"

//Verifier interface. Each verifier must implement this interface. Please do not modify it
type ValidatorInterface interface {
	CheckParams(context *gin.Context)
}
