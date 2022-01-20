package factory

import (
	"edgelog/app/core/container"
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/http/validator/core/interf"
	"github.com/gin-gonic/gin"
)

//Form parameter validator factory (do not modify)
func Create(key string) func(context *gin.Context) {

	if value := container.CreateContainersFactory().Get(key); value != nil {
		if val, isOk := value.(interf.ValidatorInterface); isOk {
			return val.CheckParams
		}
	}
	variable.ZapLog.Error(my_errors.ErrorsValidatorNotExists + ",  Verifier module ：" + key)
	return nil
}
