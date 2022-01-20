package chaptcha

import (
	"bytes"
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/utils/response"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"time"
)

type Captcha struct {
	Id      string `json:"id"`
	ImgUrl  string `json:"img_url"`
	Refresh string `json:"refresh"`
	Verify  string `json:"verify"`
}

//  Generate verification code ID
func (c *Captcha) GenerateId(context *gin.Context) {
	//  Sets the numeric length of the verification code （ number ）
	var length = variable.ConfigYml.GetInt("Captcha.length")
	captchaId := captcha.NewLen(length)
	c.Id = captchaId
	c.ImgUrl = "/captcha/" + captchaId + ".png"
	c.Refresh = c.ImgUrl + "?reload=1"
	c.Verify = "/captcha/" + captchaId + "/ Replace here with the correct verification code for verification "
	response.Success(context, " Verification code information ", c)
}

//  Get verification code image 
func (c *Captcha) GetImg(context *gin.Context) {
	captchaId := context.Param("captchaId")
	_, file := path.Split(context.Request.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || captchaId == "" {
		response.Fail(context, consts.CaptchaGetParamsInvalidCode, consts.CaptchaGetParamsInvalidMsg, "")
		return
	}

	if context.Query("reload") != "" {
		captcha.Reload(id)
	}

	context.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	context.Header("Pragma", "no-cache")
	context.Header("Expires", "0")

	var vBytes bytes.Buffer
	if ext == ".png" {
		context.Header("Content-Type", "image/png")
		//  Set the picture size of the verification code required by the actual business （ wide  X  high ），captcha.StdWidth, captcha.StdHeight  Is the default ， Please modify it to a specific number 
		_ = captcha.WriteImage(&vBytes, id, captcha.StdWidth, captcha.StdHeight)
		http.ServeContent(context.Writer, context.Request, id+ext, time.Time{}, bytes.NewReader(vBytes.Bytes()))
	}
}

//  Verification code 
func (c *Captcha) CheckCode(context *gin.Context) {
	captchaIdKey := variable.ConfigYml.GetString("Captcha.captchaId")
	captchaValueKey := variable.ConfigYml.GetString("Captcha.captchaValue")

	captchaId := context.Param(captchaIdKey)
	value := context.Param(captchaValueKey)

	if captchaId == "" || value == "" {
		response.Fail(context, consts.CaptchaCheckParamsInvalidCode, consts.CaptchaCheckParamsInvalidMsg, "")
		return
	}
	if captcha.VerifyString(captchaId, value) {
		response.Success(context, consts.CaptchaCheckOkMsg, "")
	} else {
		response.Fail(context, consts.CaptchaCheckFailCode, consts.CaptchaCheckFailMsg, "")
	}
}
