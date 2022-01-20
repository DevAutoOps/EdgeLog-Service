package users

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"github.com/gin-gonic/gin"
)

//  simulation Aop  Implement the pre and post callback of a controller function 

type DestroyBefore struct{}

//  The pre function must have a return value ， In this way, you can control whether the process continues downward 
func (d *DestroyBefore) Before(context *gin.Context) bool {
	userId := context.GetFloat64(consts.ValidatorPrefix + "id")
	variable.ZapLog.Sugar().Infof(" simulation  Users  Delete operation ， Before  Callback , user ID：%.f\n", userId)
	if userId > 10 {
		return true
	} else {
		return false
	}
}
