package home

import (
	"edgelog/app/global/consts"
	"edgelog/app/http/controller/api"
	common_data_type "edgelog/app/http/validator/common/data_type"
	"edgelog/app/http/validator/core/data_transfer"
	"edgelog/app/utils/response"
	"github.com/gin-gonic/gin"
)

//The front-end interface of portal class simulates a parameter validator for obtaining news

type News struct {
	NewsType string `form:"newsType" json:"newsType"  binding:"required,min=1"` //Validation rule: required, minimum length is 1
	common_data_type.Page
}

func (n News) CheckParams(context *gin.Context) {
	//1. According to the basic syntax provided by the verifier, more than 90% of the unqualified parameters can be verified
	if err := context.ShouldBind(&n); err != nil {
		response.ErrorParam(context, gin.H{
			"tips": "HomeNews Parameter verification failed ， The parameters do not meet the requirements ，newsType( length >=1)、page>=1、limit>=1, Please check yourself according to the rules ",
			"err":  err.Error(),
		})
		return
	}

	//This function is mainly used to directly pass the bound data to the next step (controller) in the form of key = > value
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "HomeNews Form validator json Chemical failure ", "")
	} else {
		//After verification, call the controller and pass the verifier member (field) to the controller to maintain the consistency of context data
		(&api.Home{}).News(extraAddBindDataContext)
	}

}
