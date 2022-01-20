package response

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/my_errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnJson(Context *gin.Context, httpCode int, dataCode int, msg string, data interface{}) {

	//Context.Header("key2020","value2020")  	// Additional information can be added to the header according to the actual situation 
	Context.JSON(httpCode, gin.H{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

//  take json Character channeling to standard json Format return （ for example ， from redis read json、 Formatted string ， Return to browser json format ）
func ReturnJsonFromString(Context *gin.Context, httpCode int, jsonStr string) {
	Context.Header("Content-Type", "application/json; charset=utf-8")
	Context.String(httpCode, jsonStr)
}

//  Syntax sugar function encapsulation 

//  Direct return success 
func Success(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusOK, consts.CurdStatusOkCode, msg, data)
}

func Error(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  err.Error(),
	})
	c.Abort()
}

//  Failed business logic 
func Fail(c *gin.Context, dataCode int, msg string, data interface{}) {
	ReturnJson(c, http.StatusBadRequest, dataCode, msg, data)
	c.Abort()
}

//token  Permission verification failed 
func ErrorTokenAuthFail(c *gin.Context) {
	ReturnJson(c, http.StatusUnauthorized, http.StatusUnauthorized, my_errors.ErrorsNoAuthorization, "")
	// Terminate the execution of other callback functions that may have been loaded 
	c.Abort()
}

// casbin  Authentication failed ， return  405  Method does not allow access 
func ErrorCasbinAuthFail(c *gin.Context, msg interface{}) {
	ReturnJson(c, http.StatusMethodNotAllowed, http.StatusMethodNotAllowed, my_errors.ErrorsCasbinNoAuthorization, msg)
	c.Abort()
}

// Parameter verification error 
func ErrorParam(c *gin.Context, wrongParam interface{}) {
	ReturnJson(c, http.StatusBadRequest, consts.ValidatorParamsCheckFailCode, consts.ValidatorParamsCheckFailMsg, wrongParam)
	c.Abort()
}

//  System execution code error 
func ErrorSystem(c *gin.Context, msg string, data interface{}) {
	ReturnJson(c, http.StatusInternalServerError, consts.ServerOccurredErrorCode, consts.ServerOccurredErrorMsg+msg, data)
	c.Abort()
}

type PageRes struct {
	List      interface{} `json:"rows"`
	Count     int         `json:"total"`
	PageIndex int         `json:"pageNum"`
	PageSize  int         `json:"pageSize"`
}

//Compatible function
func Custum(c *gin.Context, data gin.H) {
	c.AbortWithStatusJSON(http.StatusOK, data)
}
