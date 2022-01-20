package upload_file

import (
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"edgelog/app/utils/md5_encrypt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strings"
	"time"
)

func Upload(context *gin.Context, savePath string) (r bool, finnalSavePath interface{}) {

	newSavePath, newReturnPath := generateYearMonthPath(savePath)

	//  1. Get the uploaded file name ( The parameter validator has verified the first step error ， Simplify here )
	file, _ := context.FormFile(variable.ConfigYml.GetString("FileUploadSetting.UploadFileField")) //  file  Is a file structure （ File object ）

	//   Save file ， The original file name is globally unique encoded and encrypted 、md5  encryption ， Ensure no duplicate storage in the background 
	var saveErr error
	if sequence := variable.SnowFlake.GetId(); sequence > 0 {
		saveFileName := fmt.Sprintf("%d%s", sequence, file.Filename)
		saveFileName = md5_encrypt.MD5(saveFileName) + path.Ext(saveFileName)

		if saveErr = context.SaveUploadedFile(file, newSavePath+saveFileName); saveErr == nil {
			//   Upload succeeded , Returns the relative path of the resource ， Here, please return the absolute path or relative path according to the actual situation 
			finnalSavePath = gin.H{
				"path": strings.ReplaceAll(newReturnPath+saveFileName, variable.BasePath, ""),
			}
			return true, finnalSavePath
		}
	} else {
		saveErr = errors.New(my_errors.ErrorsSnowflakeGetIdFail)
		variable.ZapLog.Error(" Error saving file ：" + saveErr.Error())
	}
	return false, nil

}

//  File upload can be set according to  xxx year -xx month   Format storage 
func generateYearMonthPath(savePathPre string) (string, string) {
	returnPath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileReturnPath")
	curYearMonth := time.Now().Format("2006_01")
	newSavePathPre := savePathPre + curYearMonth
	newReturnPathPre := returnPath + curYearMonth
	//  The correlation path does not exist ， Create directory 
	if _, err := os.Stat(newSavePathPre); err != nil {
		if err = os.MkdirAll(newSavePathPre, os.ModePerm); err != nil {
			variable.ZapLog.Error(" Error in file upload and directory creation " + err.Error())
			return "", ""
		}
	}
	return newSavePathPre + "/", newReturnPathPre + "/"
}
