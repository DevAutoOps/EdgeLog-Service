package web

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/service/upload_file"
	"edgelog/app/utils/response"
	"github.com/gin-gonic/gin"
)

type Upload struct {
}

//   File upload is a separate module ， Return the storage path of the uploaded file to any business 。
//  Start uploading 
func (u *Upload) StartUpload(context *gin.Context) {
	savePath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath")
	if r, finnalSavePath := upload_file.Upload(context, savePath); r == true {
		response.Success(context, consts.CurdStatusOkMsg, finnalSavePath)
	} else {
		response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
	}
}
