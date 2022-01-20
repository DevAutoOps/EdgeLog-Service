package handle

import (
	"edgelog/app/global/consts"
	"edgelog/app/global/variable"
	"edgelog/app/model"
	"edgelog/app/service/remote"
	"edgelog/app/utils/response"
	"edgelog/routers/common"
	"edgelog/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
)

func HandleNode(group *gin.RouterGroup) {
	group.Use(common.TokenMid)
	group.GET("/getNode", getNode)
	group.POST("/editNode", editNode)
	group.POST("/initNode", initNode)
}

// @Summary  Query agent status
// @Tags  Node
// @Success 200 {object} []model.Node "{"code": 200, "data": [...]}"
// @Router /api/v1/node/getNode [get]
func getNode(context *gin.Context) {
	response.Success(context, consts.CurdStatusOkMsg, []model.Node{*variable.Node})
}

// @Summary  Query agent status
// @Tags  Node
// @Param name body string true "node name"
// @Param ip body string true "node ip"
// @Param port body int true "node port"
// @Param account body string true "ssh account"
// @Param password body string true "ssh password"
// @Param agentPort body int true "agent port"
// @Param os body int true "system type"
// @Param conf body string true "nginx conf path"
// @Param logs body string true "nginx logs path"
// @Param templateId body int true "template id"
// @Success 200 {object} model.Node "{"code": 200, "data": [...]}"
// @Router /api/v1/node/editNode [post]
func editNode(context *gin.Context) {
	node := &model.Node{}
	if err := context.ShouldBind(node); err != nil {
		common.Error(context, err)
		return
	}
	if len(node.Name) > 0 {
		variable.Node.Name = node.Name
	}
	if len(node.Ip) > 0 {
		variable.Node.Ip = node.Ip
	}
	if node.Port > 0 {
		variable.Node.Port = node.Port
	}
	if len(node.Account) > 0 {
		variable.Node.Account = node.Account
	}
	if len(node.Password) > 0 {
		variable.Node.Password = node.Password
	}
	if node.AgentPort > 0 {
		variable.Node.AgentPort = node.AgentPort
	}
	variable.Node.Os = node.Os
	if len(node.Conf) > 0 {
		variable.Node.Conf = node.Conf
	}
	if len(node.Logs) > 0 {
		variable.Node.Logs = node.Logs
	}
	if node.TemplateId > 0 {
		variable.Node.TemplateId = node.TemplateId
	}
	response.Success(context, consts.CurdStatusOkMsg, variable.Node)
}

// @Summary  node agent install
// @Tags  base / node
// @Router /api/v1/node/initNode [post]
func initNode(context *gin.Context) {
	ip, err := tools.ExternalIP()
	if err != nil {
		variable.ZapLog.Error(" Probe initialization failed ,"+consts.IpNotFoundMsg, zap.Any("err", err))
		response.Custum(context, gin.H{
			"msg":  " Probe initialization failed ," + consts.IpNotFoundMsg,
			"code": consts.IpNotFoundCode,
			"data": fmt.Sprintf("%v", err),
		})
		return
	}

	remoteCode, remoteMsg := nodeInitSet(ip.String())
	if remoteCode == 0 {
		variable.Node.IsInit = true
		response.Success(context, consts.CurdStatusOkMsg, remoteMsg)
	} else {
		returnCode := remoteCode
		if remoteCode%2 == 0 {
			returnCode = 2
			remoteMsg += consts.InitPermissionDeniedMsg
		}
		variable.ZapLog.Error(" Probe initialization failed ", zap.Any("err", err))
		response.Custum(context, gin.H{
			"msg":  " Probe initialization failed ",
			"code": returnCode,
			"data": remoteMsg,
		})
	}
	variable.ZapLog.Info(fmt.Sprintf(" Server node probe initialization log ， Return code ：%d, msg: %s", remoteCode, remoteMsg))
}

func nodeInitSet(ip string) (int, string) {
	remoteCode := 0
	service := remote.RemoteService{}
	remoteMsg, err := service.Connect(variable.Node.Account, variable.Node.Password, "", variable.Node.Ip, tools.IntToString(variable.Node.Port), 0)
	if err != nil {
		remoteMsg += " This account cannot log in to the server ， Please check whether the account and password are correct \n"
		remoteCode = -1
		return remoteCode, remoteMsg
	}
	defer service.Release()

	install := "/usr/local"
	agentName := "agent"
	uploadPath := filepath.Join(variable.BasePath, "/public/install/agent.tar.gz")
	remoteAgentDir := strings.Replace(filepath.Join(install, agentName), "\\", "/", -1)
	remoteAgentFile := agentName

	//  Probe initialization
	//  Upload file
	cmdRmAgentPackage := fmt.Sprintf("rm -rf %s", strings.Replace(filepath.Join(install, filepath.Base(uploadPath)), "\\", "/", -1))
	_, _, _ = service.CommandSingleWithErr(cmdRmAgentPackage)

	putMsg, err := service.PutOne(uploadPath, install+"/")
	remoteMsg += putMsg
	if err != nil {
		fmt.Printf(" Upload file %s fail ： %s", uploadPath, err.Error())
		if strings.Contains(err.Error(), consts.PermissionDenied) {
			remoteCode = 2
		} else {
			remoteCode = 3
		}
		//  Return after failed to upload file
		return remoteCode, remoteMsg
	}

	//  Execute predetermined instructions
	cmdOrder := "cd /usr/local/;tar -zxf agent.tar.gz"
	if len(cmdOrder) > 0 {
		cmdOrder = strings.Replace(cmdOrder, ";", "&&", -1)
		output2, output2Msg, err := service.CommandSingleWithErr(cmdOrder)
		remoteMsg += output2
		remoteMsg += output2Msg
		if len(output2) > 0 && strings.Contains(output2, consts.PermissionDenied) {
			remoteCode = 4
			return remoteCode, remoteMsg
		} else if err != nil {
			remoteCode = 5
			return remoteCode, remoteMsg
		}
	}

	//  Read probe configuration file
	remoteConfigDir := strings.Replace(filepath.Join(remoteAgentDir, "config"), "\\", "/", -1)
	remoteAgentConfigName := consts.AgentConfigName
	remoteConfigYml := ""
	//  Read configuration file
	cmdCatConfig := fmt.Sprintf("cat %s", strings.Replace(filepath.Join(remoteConfigDir, remoteAgentConfigName), "\\", "/", -1))
	outputCatConfig, outputCatConfigMsg, err := service.CommandSingleWithErr(cmdCatConfig)
	remoteMsg += outputCatConfigMsg
	if len(outputCatConfig) > 0 && strings.Contains(outputCatConfig, consts.PermissionDenied) {
		remoteCode = 6
		return remoteCode, remoteMsg
	} else if err != nil || len(outputCatConfig) == 0 {
		remoteCode = 7
		return remoteCode, remoteMsg
	}
	remoteConfigYml = outputCatConfig

	//  Modify profile
	agentConfig, err := tools.ReadYaml(remoteConfigYml)
	if err != nil {
		remoteCode = 7
		remoteMsg += fmt.Sprintf(" serialize agent Error in configuration file , err: %s\n", err)
		return remoteCode, remoteMsg
	}
	agentConfig.Server.Ip = ip
	agentConfig.Server.MonitorPort = fmt.Sprintf(":%d", variable.ConfigYml.GetInt("HttpServer.ProbesMonitorPort"))
	agentConfig.Server.LogPort = fmt.Sprintf(":%d", variable.ConfigYml.GetInt("HttpServer.ProbesLogPort"))
	agentConfig.Modules.LogCollector.Enable = true
	agentConfig.Modules.LogCollector.Targets = []string{"Nginx"}
	agentConfig.Modules.LogCollector.TargetConfs = []string{variable.Node.Conf}
	agentConfig.Modules.Monitor.Enable = true
	agentConfig.Modules.Monitor.Freq = 5000
	agentConfig.Modules.Monitor.Check = true
	agentConfig.Modules.Monitor.CheckFreq = 2
	agentConfigResult, err := tools.WriteYaml(agentConfig)
	if err != nil {
		remoteCode = 9
		remoteMsg += fmt.Sprintf(" Deserialization agent Error in configuration file , err: %s\n", err)
		return remoteCode, remoteMsg
	}
	if agentConfigResult == "" {
		remoteCode = 9
		remoteMsg += fmt.Sprintf(" Deserialization agent The configuration file is empty , err: %s\n", err)
		return remoteCode, remoteMsg
	}
	cmdSaveConfigXml := "cat << EOF > " + remoteConfigDir + "/" + remoteAgentConfigName + "\n" + agentConfigResult + "\nEOF\n\n"
	outputSaveConfigXml, outputSaveConfigXmlMsg, err := service.CommandSingleWithErr(cmdSaveConfigXml)
	remoteMsg += outputSaveConfigXmlMsg
	if len(outputSaveConfigXml) > 0 && strings.Contains(outputSaveConfigXml, consts.PermissionDenied) {
		remoteCode = 8
		return remoteCode, remoteMsg
	} else if err != nil {
		remoteCode = 9
		return remoteCode, remoteMsg
	}

	//  Add executable permissions to the probe
	cmdAddExe := fmt.Sprintf("chmod +x %s", strings.Replace(filepath.Join(remoteAgentDir, remoteAgentFile), "\\", "/", -1))
	outputAddExe, outputAddExeMsg, err := service.CommandSingleWithErr(cmdAddExe)
	remoteMsg += outputAddExe
	remoteMsg += outputAddExeMsg
	if len(outputAddExeMsg) > 0 && strings.Contains(outputAddExeMsg, consts.PermissionDenied) {
		remoteCode = 10
		return remoteCode, remoteMsg
	} else if err != nil {
		remoteCode = 11
		return remoteCode, remoteMsg
	}

	//  View process kill
	cmdKill := fmt.Sprintf("kill -s 9 `ps -aux | grep %s | grep -v grep | awk '{print $2}'`", remoteAgentFile)
	outputKill, outputKillMsg, err := service.CommandSingleWithErr(cmdKill)
	remoteMsg += outputKill
	remoteMsg += outputKillMsg
	if len(outputKill) > 0 && strings.Contains(outputKill, consts.PermissionDenied) {
		remoteCode = 12
		return remoteCode, remoteMsg
	}

	//  start-up
	cmdStart := fmt.Sprintf("cd %s;nohup ./%s > log.txt 2>&1 &", remoteAgentDir, remoteAgentFile)
	outputStart, outputStartMsg, err := service.CommandSingleWithErr(cmdStart)
	remoteMsg += outputStart
	remoteMsg += outputStartMsg
	if len(outputStart) > 0 && strings.Contains(outputStart, consts.PermissionDenied) {
		remoteCode = 14
		return remoteCode, remoteMsg
	} else if err != nil {
		remoteCode = 15
		return remoteCode, remoteMsg
	}

	//  delete tar
	cmdRmAgentPackage = fmt.Sprintf("rm -rf %s", strings.Replace(filepath.Join(install, filepath.Base(uploadPath)), "\\", "/", -1))
	_, _, _ = service.CommandSingleWithErr(cmdRmAgentPackage)

	return remoteCode, remoteMsg
}
