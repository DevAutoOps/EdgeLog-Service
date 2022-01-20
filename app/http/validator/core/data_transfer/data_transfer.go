package data_transfer

import (
	"edgelog/app/global/variable"
	"edgelog/app/http/validator/core/interf"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

//The verifier member (field) is bound to the data transmission context, which is convenient for the controller to obtain
/**
Parameter description of this function:
Validatorinterface implements the structure of the validator interface
extra_ add_ data_ Prefix the data prefix passed to the controller by the validator binding parameter
Context gin context
*/

func DataAddContext(validatorInterface interf.ValidatorInterface, extraAddDataPrefix string, context *gin.Context) *gin.Context {
	var tempJson interface{}
	if tmpBytes, err1 := json.Marshal(validatorInterface); err1 == nil {
		if err2 := json.Unmarshal(tmpBytes, &tempJson); err2 == nil {
			if value, ok := tempJson.(map[string]interface{}); ok {
				for k, v := range value {
					context.Set(extraAddDataPrefix+k, v)
				}
				//In addition, three keys are appended to the context: created_ at  、 updated_ at  、 deleted_ At, select and obtain relevant key values according to actual needs
				curDateTime := time.Now().Format(variable.DateFormat)
				context.Set(extraAddDataPrefix+"created_at", curDateTime)
				context.Set(extraAddDataPrefix+"updated_at", curDateTime)
				context.Set(extraAddDataPrefix+"deleted_at", curDateTime)
				return context
			}
		}
	}
	return nil
}
