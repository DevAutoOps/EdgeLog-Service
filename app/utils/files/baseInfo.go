package files

import (
	"edgelog/app/global/my_errors"
	"edgelog/app/global/variable"
	"mime/multipart"
	"net/http"
	"os"
)

//  Return value description ：
//	7z、exe、doc  Type will return  application/octet-stream   Unknown file type 
//	jpg	=>	image/jpeg
//	png	=>	image/png
//	ico	=>	image/x-icon
//	bmp	=>	image/bmp
//  xlsx、docx 、zip	=>	application/zip
//  tar.gz	=>	application/x-gzip
//  txt、json、log Text files such as 	=>	text/plain; charset=utf-8    remarks ： even if txt yes gbk、ansi code ， Will also be identified as utf-8

//  Get file by file name mime information 
func GetFilesMimeByFileName(filepath string) string {
	f, err := os.Open(filepath)
	if err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadOpenFail + err.Error())
	}
	defer f.Close()

	//  Just before  32  Just a byte 
	buffer := make([]byte, 32)
	if _, err := f.Read(buffer); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadReadFail + err.Error())
		return ""
	}

	return http.DetectContentType(buffer)
}

//  Get file through file pointer mime information 
func GetFilesMimeByFp(fp multipart.File) string {

	buffer := make([]byte, 32)
	if _, err := fp.Read(buffer); err != nil {
		variable.ZapLog.Error(my_errors.ErrorsFilesUploadReadFail + err.Error())
		return ""
	}

	return http.DetectContentType(buffer)
}
