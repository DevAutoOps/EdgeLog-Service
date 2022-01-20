package upload_files

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/http/controller/web"
	"edgelog/app/utils/files"
	"edgelog/app/utils/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type UpFiles struct {
}

//File upload common module form parameter verifier
func (u UpFiles) CheckParams(context *gin.Context) {
	tmpFile, err := context.FormFile(variable.ConfigYml.GetString("FileUploadSetting.UploadFileField")) //File is a file structure (file object)
	var isPass bool
	//There is an error in obtaining the file. An empty file may be uploaded
	if err != nil {
		response.Fail(context, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, err.Error())
		return
	}
	//If it exceeds the maximum value set by the system: 32m, the unit of tmpfile.size is bytes. Compared with the file unit m defined by us, we need to compare our unit * 1024 * 1024 (i.e. the 20th power of 2), which is < < 20 in one step
	if tmpFile.Size > variable.ConfigYml.GetInt64("FileUploadSetting.Size")<<20 {
		response.Fail(context, consts.FilesUploadMoreThanMaxSizeCode, consts.FilesUploadMoreThanMaxSizeMsg+variable.ConfigYml.GetString("FileUploadSetting.Size"), "")
		return
	}
	//File MIME type not allowed
	if fp, err := tmpFile.Open(); err == nil {
		mimeType := files.GetFilesMimeByFp(fp)

		for _, value := range variable.ConfigYml.GetStringSlice("FileUploadSetting.AllowMimeType") {
			if strings.ReplaceAll(value, " ", "") == strings.ReplaceAll(mimeType, " ", "") {
				isPass = true
				break
			}
		}
		_ = fp.Close()
	} else {
		response.ErrorSystem(context, consts.ServerOccurredErrorMsg, "")
		return
	}
	//If there are equal types, call the controller through verification
	if !isPass {
		response.Fail(context, consts.FilesUploadMimeTypeFailCode, consts.FilesUploadMimeTypeFailMsg, "")
	} else {
		(&web.Upload{}).StartUpload(context)
	}
}
