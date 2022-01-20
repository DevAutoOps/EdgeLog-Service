package tools

import (
	"edgelog/app/global/consts"
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func GetOs() string {
	return strings.ToLower(runtime.GOOS)
}

func CmdLinux(cmdStr string) (result string, success bool, err error) {
	result = ""
	success = false
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		//Setpgid: true,	//  Comment out this line ，windows This structure does not have this attribute under compilation
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result = string(out)
	success = true
	return
}

func CmdWindows(cmdStr string) (result string, success bool, err error) {
	result = ""
	success = false
	cmd := exec.Command("cmd", "/C", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	result = string(out)
	success = true
	return
}

func Cmd(cmdStr string) (result string, success bool, err error) {
	switch GetOs() {
	case consts.OsWindows:
		return CmdWindows(cmdStr)
	case consts.OsLinux:
		return CmdLinux(cmdStr)
	default:
		return "", false, errors.New("I won't support it")
	}
}

func CmdSerialize(cmdStr string) string {
	return strings.Replace(cmdStr, " ", "_", -1)
}

func CmdDeSerialize(cmdStr string) string {
	return strings.Replace(cmdStr, "_", " ", -1)
}
