package remote

import (
	"edgelog/app/global/variable"
	ssh2 "edgelog/app/utils/ssh"
	"fmt"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
)

func init() {
	cliMutex = &sync.Mutex{}
}

type Client struct {
	Cli      *ssh2.Client
	HomePath string
}

var (
	cliMutex *sync.Mutex
	cliMap   = map[string]*Client{}
	wg       sync.WaitGroup
)

//Start a parallel task
func launch(f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}

//Done built-in command, waiting for parallel tasks to complete
func Done() {
	wg.Wait()
	log.Println("multi command done")
}

//Connect built-in command to connect to the server
func Connect(user, password, keyFile, host, port string, types int) (msg string, err error) {
	//launch(func() {
	wg.Add(1)
	go func(user, password, keyFile, host, port string, types int) {
		defer wg.Done()
		//In order to prevent the connection from taking too long, another thread is opened to handle it
		remoteAdd := host + ":" + port
		if _, ok := cliMap[remoteAdd]; ok {
			//The connected server will not be connected again
			log.Printf("[%s] connected\n\n", remoteAdd)
			msg += fmt.Sprintf("[%s] connected\n\n", remoteAdd)
			return
		}

		if types >= 2 {
			log.Println("Currently only password and key connections are supported")
			msg += fmt.Sprintf("Currently only password and key connections are supported\n")
			return
		}
		var client *ssh2.Client
		if types == 0 {
			client, err = ssh2.NewClient(host, port, user, password)
		} else {
			client, err = ssh2.NewWithPrivateKey(host, port, user, keyFile)
		}

		if err != nil {
			log.Printf("[%s] ssh Dial error: %s\n", remoteAdd, err)
			msg += fmt.Sprintf("[%s] ssh Dial error: %s\n", remoteAdd, err)
			variable.ZapLog.Error("[sshClient] ssh Dial error:", zap.String("host:port", remoteAdd), zap.Error(err))
			return
		}

		homePath, err := client.Output("pwd")
		if err != nil {
			log.Printf("[%s] get home path error: %s", remoteAdd, err.Error())
			msg += fmt.Sprintf("[%s] get home path error: %s", remoteAdd, err.Error())
			variable.ZapLog.Error("[sshClient] get home path error:", zap.String("host:port", remoteAdd), zap.Error(err))
			return
		}

		cli := &Client{
			Cli:      client,
			HomePath: strings.TrimSpace(string(homePath)),
		}

		cliMutex.Lock()
		defer cliMutex.Unlock()
		cliMap[remoteAdd] = cli
		log.Printf("[%s] connect success\n", remoteAdd)
		msg += fmt.Sprintf("[%s] connect success\n", remoteAdd)
	}(user, password, keyFile, host, port, types)
	Done()
	//})
	return
}

//Release built-in command to release the connection
func Release(host, port string) {
	remoteAdd := host + ":" + port
	if client, ok := cliMap[remoteAdd]; ok {
		err := client.Cli.Close()
		if err != nil {
			log.Printf("[%s] close error: %s\n", remoteAdd, err)
			variable.ZapLog.Error("[sshClient] close error:", zap.String("host:port", remoteAdd), zap.Error(err))
		}
		delete(cliMap, remoteAdd)
		log.Printf("[%s] released\n", remoteAdd)
		return
	}
	log.Printf("[%s] has not connected yet\n", remoteAdd)
}

func ReleaseOne(host string) {
	if client, ok := cliMap[host]; ok {
		err := client.Cli.Close()
		if err != nil {
			log.Printf("[%s] close error: %s\n", host, err)
			variable.ZapLog.Error("[sshClient] close error:", zap.String("host:port", host), zap.Error(err))
		}
		delete(cliMap, host)
		log.Printf("[%s] released\n", host)
		return
	}
	log.Printf("[%s] has not connected yet\n", host)
}

//The ReleaseAll built-in command releases all connections
func ReleaseAll() {
	for host := range cliMap {
		if client, ok := cliMap[host]; ok {
			err := client.Cli.Close()
			if err != nil {
				log.Printf("[%s] close error: %s\n", host, err)
				variable.ZapLog.Error("[sshClient] close error:", zap.String("host:port", host), zap.Error(err))
			}
			//delete(cliMap, host)
			log.Printf("[%s] released\n", host)
			continue
		}
		log.Printf("[%s] has not connected yet\n", host)
	}
	cliMap = map[string]*Client{}
}

//Remote built-in command, batch remote execution
func Remote(cmd string) {
	for host, client := range cliMap {
		fmt.Printf("\033[36m>>>>>>>>>>>>>>> %s [%s] <<<<<<<<<<<<<<<\033[0m\n", host, cmd)
		err := client.Cli.Exec(cmd)
		if err != nil {
			//There is an error in this execution
			log.Printf("[%s] remote command [%s] failed: %s\n", host, cmd, err.Error())
			variable.ZapLog.Error("[sshClient] remote command error:", zap.String("host:port", host), zap.String("cmd", cmd), zap.Error(err))
			continue
		}
		log.Printf("[%s] remote command [%s] success\n\n", host, cmd)
	}
}

//Execute a single server command remotely and return a message
func RemoteSingle(cmd string) string {
	var output []byte
	for host, client := range cliMap {
		fmt.Printf("\033[36m>>>>>>>>>>>>>>> %s [%s] <<<<<<<<<<<<<<<\033[0m\n", host, cmd)
		var err error
		output, err = client.Cli.Output(cmd)
		if err != nil {
			//There is an error in this execution
			log.Printf("[%s] remote command [%s] failed: %s\n", host, cmd, err.Error())
			variable.ZapLog.Error("[sshClient] remote command failed:", zap.String("host:port", host), zap.String("cmd", cmd), zap.Error(err))
			continue
		}
		log.Printf("[%s] remote command [%s] success\n\n", host, cmd)
	}
	return string(output)
}

//Remotely execute server commands and return standard output and error output messages
func RemoteWithAllOutput(cmd string) (string, string, error) {
	output2 := ""
	outputStr := ""
	var err error
	for host, client := range cliMap {
		fmt.Printf("\033[36m>>>>>>>>>>>>>>> %s [%s] <<<<<<<<<<<<<<<\033[0m\n", host, cmd)
		var output []byte
		output, err = client.Cli.OutputAll(cmd)
		output2 += string(output)
		if err != nil {
			//There is an error in this execution
			log.Printf("[%s] remote command [%s] failed: %s\n", host, cmd, err.Error())
			variable.ZapLog.Error("[sshClient] remote command failed:", zap.String("host:port", host), zap.String("cmd", cmd), zap.Error(err))
			continue
		}
		log.Printf("[%s] remote command [%s] success\n\n", host, cmd)
		outputStr += fmt.Sprintf("[%s] remote command [%s] success\n\n", host, cmd)
	}
	return output2, outputStr, err
}

//Remotely execute server commands and return standard output and error output messages
func RemoteWithAllOutput2(cmd string) (string, string, string, error) {
	output2 := ""
	outputStr := ""
	output3 := ""
	var err error
	for host, client := range cliMap {
		fmt.Printf("\033[36m>>>>>>>>>>>>>>> %s [%s] <<<<<<<<<<<<<<<\033[0m\n", host, cmd)
		var output []byte
		output, err = client.Cli.OutputAll(cmd)
		output2 += fmt.Sprintf("[%s]: %s", host, string(output))
		if err != nil {
			//There is an error in this execution
			log.Printf("[%s] remote command [%s] failed: %s\n", host, cmd, err.Error())
			output3 += fmt.Sprintf("[%s] remote command [%s] failed: %s\n", host, cmd, string(output))
			variable.ZapLog.Error("[sshClient] remote command failed:", zap.String("host:port", host), zap.String("cmd", cmd), zap.Error(err))
			continue
		}
		log.Printf("[%s] remote command [%s] success\n\n", host, cmd)
		outputStr += fmt.Sprintf("[%s] remote command [%s] success\n\n", host, cmd)
		output3 += outputStr + output2 + "\n"
	}
	return output2, outputStr, output3, err
}

//Remotely execute commands on a single server and return standard output and error output messages
func RemoteWithAllOutput3(cmd string, host string) (string, string, string, error) {
	output2 := ""
	outputStr := ""
	output3 := ""
	var err error
	if client, ok := cliMap[host]; ok {
		fmt.Printf("\033[36m>>>>>>>>>>>>>>> %s [%s] <<<<<<<<<<<<<<<\033[0m\n", host, cmd)
		var output []byte
		output, err = client.Cli.OutputAll(cmd)
		output2 += fmt.Sprintf("[%s]: %s", host, string(output))
		if err != nil {
			//There is an error in this execution
			log.Printf("[%s] remote command [%s] failed: %s\n", host, cmd, err.Error())
			output3 += fmt.Sprintf("[%s] remote command [%s] failed: %s\n", host, cmd, string(output))
			variable.ZapLog.Error("[sshClient] remote command failed:", zap.String("host:port", host), zap.String("cmd", cmd), zap.Error(err))

		} else {
			log.Printf("[%s] remote command [%s] success\n\n", host, cmd)
			outputStr += fmt.Sprintf("[%s] remote command [%s] success\n\n", host, cmd)
			output3 += outputStr + output2 + "\n"
		}
	}
	return output2, outputStr, output3, err
}

//Remoteonewithalloutput remotely executes single server commands and returns standard output and error output messages
func RemoteOneWithAllOutput(host, cmd string) (string, string, error) {
	output2 := ""
	outputStr := ""
	var err error
	if client, ok := cliMap[host]; ok {
		fmt.Printf("\033[36m>>>>>>>>>>>>>>> %s [%s] <<<<<<<<<<<<<<<\033[0m\n", host, cmd)
		var output []byte
		output, err = client.Cli.OutputAll(cmd)
		output2 += string(output)
		if err != nil {
			//There is an error in this execution
			log.Printf("[%s] remote command [%s] failed: %s\n", host, cmd, err.Error())
			variable.ZapLog.Error("[sshClient] remote command failed:", zap.String("host:port", host), zap.String("cmd", cmd), zap.Error(err))
		}
		log.Printf("[%s] remote command [%s] success\n\n", host, cmd)
		outputStr += fmt.Sprintf("[%s] remote command [%s] success\n\n", host, cmd)
	}
	return output2, outputStr, err
}

//Check built-in command to detect the established connection
func Check() {
	for host, client := range cliMap {
		_ = client
		log.Printf("[%s] connecting\n", host)
	}
}

//Put built-in command, batch upload files
func Put(local string, dstDir string) (msg string) {
	fmt.Println("func put:", local, dstDir)
	for host, client := range cliMap {
		wg.Add(1)
		go func(host string, client *Client, local, dstDir string) {
			defer wg.Done()
			toDir := dstDir
			if toDir == "" {
				toDir = client.HomePath
			}
			err := client.Cli.Upload(local, dstDir)
			if err != nil {
				log.Printf("[%s] upload file %s error: %s\n", host, local, err)
				msg += fmt.Sprintf("[%s] upload file %s error: %s\n", host, local, err)
				variable.ZapLog.Error("[sshClient] upload file error:", zap.String("host:port", host), zap.String("localFile", local), zap.Error(err))
				return
			}
			log.Printf("[%s] put file [%s] to [%s] success\n", host, local, dstDir)
			msg += fmt.Sprintf("[%s] put file [%s] to [%s] success\n", host, local, dstDir)
		}(host, client, local, dstDir)
	}
	Done()
	return
}

//Put upload files from a single remote server
func PutOne(local string, dstDir string) (msg string, err error) {
	fmt.Println("func put:", local, dstDir)
	for host, client := range cliMap {
		wg.Add(1)
		go func(host string, client *Client, local, dstDir string) {
			defer wg.Done()
			toDir := dstDir
			if toDir == "" {
				toDir = client.HomePath
			}
			err = client.Cli.Upload(local, dstDir)
			if err != nil {
				log.Printf("[%s] upload file %s error: %s\n", host, local, err)
				msg += fmt.Sprintf("[%s] upload file %s error: %s\n", host, local, err)
				variable.ZapLog.Error("[sshClient] upload file error:", zap.String("host:port", host), zap.String("localFile", local), zap.Error(err))
				return
			}
			log.Printf("[%s] put file [%s] to [%s] success\n", host, local, dstDir)
			msg += fmt.Sprintf("[%s] put file [%s] to [%s] success\n", host, local, dstDir)
		}(host, client, local, dstDir)
	}
	Done()
	return
}

//Putonev2 uploading files from a single remote server
func PutOneV2(host, local, dstDir string) (msg string, err error) {
	fmt.Println("func put:", local, dstDir)
	if client, ok := cliMap[host]; ok {
		wg.Add(1)
		go func(host string, client *Client, local, dstDir string) {
			defer wg.Done()
			toDir := dstDir
			if toDir == "" {
				toDir = client.HomePath
			}
			err = client.Cli.Upload(local, dstDir)
			if err != nil {
				log.Printf("[%s] upload file %s error: %s\n", host, local, err)
				msg += fmt.Sprintf("[%s] upload file %s error: %s\n", host, local, err)
				variable.ZapLog.Error("[sshClient] upload file error:", zap.String("host:port", host), zap.String("localFile", local), zap.Error(err))
				return
			}
			log.Printf("[%s] put file [%s] to [%s] success\n", host, local, dstDir)
			msg += fmt.Sprintf("[%s] put file [%s] to [%s] success\n", host, local, dstDir)
		}(host, client, local, dstDir)
	}
	Done()
	return
}

//Get built-in command, batch download files
func Get(remoteDir, local string) {
	for host, client := range cliMap {
		wg.Add(1)
		go func(host string, client *Client, remoteDir, local string) {
			defer wg.Done()
			err := client.Cli.Download(remoteDir, local)
			if err != nil {
				log.Printf("[%s] download file or dir %s to local %s error: %s\n", host, remoteDir, local, err)
				variable.ZapLog.Error("[sshClient] download file or dir error:", zap.String("host:port", host), zap.String("remoteFile", remoteDir), zap.String("localFile", local), zap.Error(err))
				return
			}
			log.Printf("[%s] put file [%s] to [%s] success\n", host, remoteDir, local)
		}(host, client, remoteDir, local)
	}
	Done()
}
