package api

import (
	"edgelog/app/global/consts"
	"edgelog/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Home struct {
}

// 1. Portal homepage news 
func (u *Home) News(context *gin.Context) {

	//   Because the skeleton of this project has added the fields of the form validator ( member ) Bind in context ， Therefore, you can  GetString()、GetInt64()、GetFloat64（） Get the required data type quickly 
	//  Of course, it can also be passed gin The context of the framework is obtained from the original method ， for example ： context.PostForm("name")  obtain ， The data format obtained in this way is text ， You need to continue the conversion yourself 
	newsType := context.GetString(consts.ValidatorPrefix + "newsType")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	userIp := context.ClientIP()

	//  Any data return is simulated here 
	response.Success(context, "ok", gin.H{
		"newsType": newsType,
		"page":     page,
		"limit":    limit,
		"userIp":   userIp,
		"title":    " Portal home page company news title 001",
		"content":  " Portal news content 001",
	})
}
