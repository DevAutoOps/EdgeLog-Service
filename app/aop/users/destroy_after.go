package users

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"github.com/gin-gonic/gin"
)

//  simulation Aop  Implement the pre and post callback of a controller function 

type DestroyAfter struct{}

func (d *DestroyAfter) After(context *gin.Context) {
	//  Post functions can be executed asynchronously 
	go func() {
		userId := context.GetFloat64(consts.ValidatorPrefix + "id")
		variable.ZapLog.Sugar().Infof(" simulation  Users  Delete operation ， After  Callback , user ID：%.f\n", userId)
	}()
}
