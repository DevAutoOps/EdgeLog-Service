package data_bind

import (
	"edgelog/app/global/consts"
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
)

const (
	modelStructMustPtr = "modelStruct  A pointer must be passed "
)

//  binding form The form validator has verified the completed parameters to  model  structural morphology ,
// mode  Structure supports anonymous nesting 
//  Data binding principles ： model  Defined structure fields and form validator structure settings json Label name 、 Consistent data type ， Can be bound 

func ShouldBindFormDataToModel(c *gin.Context, modelStruct interface{}) error {
	mTypeOf := reflect.TypeOf(modelStruct)
	if mTypeOf.Kind() != reflect.Ptr {
		return errors.New(modelStructMustPtr)
	}
	mValueOf := reflect.ValueOf(modelStruct)

	// analysis  modelStruct  field 
	mValueOfEle := mValueOf.Elem()
	mtf := mValueOf.Elem().Type()
	fieldNum := mtf.NumField()
	for i := 0; i < fieldNum; i++ {
		if !mtf.Field(i).Anonymous && mtf.Field(i).Type.Kind() != reflect.Struct {
			fieldSetValue(c, mValueOfEle, mtf, i)
		} else if mtf.Field(i).Type.Kind() == reflect.Struct {
			// Processing structure ( famous + anonymous )
			mValueOfEle.Field(i).Set(analysisAnonymousStruct(c, mValueOfEle.Field(i)))
		}
	}
	return nil
}

//  Analyze anonymous structures , And get the value of the anonymous structure 
func analysisAnonymousStruct(c *gin.Context, value reflect.Value) reflect.Value {

	typeOf := value.Type()
	fieldNum := typeOf.NumField()
	newStruct := reflect.New(typeOf)
	newStructElem := newStruct.Elem()
	for i := 0; i < fieldNum; i++ {
		fieldSetValue(c, newStructElem, typeOf, i)
	}
	return newStructElem
}

//  Assign values to structure fields 
func fieldSetValue(c *gin.Context, valueOf reflect.Value, typeOf reflect.Type, colIndex int) {
	relaKey := typeOf.Field(colIndex).Tag.Get("json")
	if relaKey != "-" {
		relaKey = consts.ValidatorPrefix + typeOf.Field(colIndex).Tag.Get("json")
		switch typeOf.Field(colIndex).Type.Kind() {
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			valueOf.Field(colIndex).SetInt(int64(c.GetFloat64(relaKey)))
		case reflect.Float32, reflect.Float64:
			valueOf.Field(colIndex).SetFloat(c.GetFloat64(relaKey))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			valueOf.Field(colIndex).SetUint(uint64(c.GetFloat64(relaKey)))
		case reflect.String:
			valueOf.Field(colIndex).SetString(c.GetString(relaKey))
		case reflect.Bool:
			valueOf.Field(colIndex).SetBool(c.GetBool(relaKey))
		default:
			// model  If there is a date time field ， Please set it as a string 
		}
	}
}
