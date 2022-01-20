package remote

import (
	"edgelog/app/global/variable"
)

type RemoteService struct {
}

func (rs *RemoteService) Connect(account string, password string, keyFile string, host string, port string, types int) (string, error) {
	keyFile = variable.BasePath + keyFile
	return Connect(account, password, keyFile, host, port, types)
}

//Batch release remote server connections
func (rs *RemoteService) Release() {
	ReleaseAll()
}

//Batch release remote server connections
func (rs *RemoteService) ReleaseOne(host string) {
	ReleaseOne(host)
}

//Batch execution of remote server commands
func (rs *RemoteService) Command(cmd string) {
	Remote(cmd)
}

//Batch execution of remote server commands with return message
func (rs *RemoteService) CommandWithErr(cmd string) (string, string, string, error) {
	return RemoteWithAllOutput2(cmd)
}

//Execute commands in a single server and return messages
func (rs *RemoteService) CommandOne(cmd string, host string) (string, string, string, error) {
	return RemoteWithAllOutput3(cmd, host)
}

//Execute a single remote server command and return a message
func (rs *RemoteService) CommandSingle(cmd string) string {
	return RemoteSingle(cmd)
}

//Remotely execute a single server command and return standard output and error output messages
func (rs *RemoteService) CommandSingleWithErr(cmd string) (string, string, error) {
	return RemoteWithAllOutput(cmd)
}

//Remotely execute a single server command and return standard output and error output messages
func (rs *RemoteService) CommandOneWithErr(host, cmd string) (string, string, error) {
	return RemoteOneWithAllOutput(host, cmd)
}

//Batch check remote server status
func (rs *RemoteService) Check() {
	Check()
}

//Batch upload files or folders to remote servers
func (rs *RemoteService) Put(local string, dstDir string) string {
	return Put(local, dstDir)
}

//Upload files or folders from a single remote server
func (rs *RemoteService) PutOne(local string, dstDir string) (string, error) {
	return PutOne(local, dstDir)
}

//Upload files or folders from a single remote server
func (rs *RemoteService) PutOneV2(host, local, dstDir string) (string, error) {
	return PutOneV2(host, local, dstDir)
}

//Batch download files or folders from remote servers
func (rs *RemoteService) Get(remoteDir, local string) {
	Get(remoteDir, local)
}
