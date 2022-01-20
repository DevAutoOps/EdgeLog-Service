package ssh

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//Upload upload local file to SFTP remote directory remote like Rsync
//rsync -av src/ dst     ./src/* --> /root/dst/*
//rsync -av src/ dst/    ./src/* --> /root/dst/*
//rsync -av src  dst     ./src/* --> /root/dst/src/*
//rsync -av src  dst/    ./src/* --> /root/dst/src/*
func (c *Client) Upload(local string, remote string) (err error) {
	//var localDir, localFile, remoteDir, remoteFile string

	info, err := os.Stat(local)
	if err != nil {
		return errors.New("sftp:  Skip upload  Upload(\"" + local + "\") , The local file does not exist or has the wrong format !")
	}
	if info.IsDir() {
		log.Println("sftp: UploadDir", local)
		return c.UploadDir(local, remote)
	}
	return c.UploadFile(local, remote)
}

//Download Download SFTP remote file remote to local like Rsync
func (c *Client) Download(remote string, local string) (err error) {
	if c.IsNotExist(strings.TrimSuffix(remote, "/")) {
		return errors.New("sftp:  The remote file does not exist , Skip file download  \"" + remote + "\" ")
	}
	if c.IsDir(remote) {
		//Return errors. New ("it is detected that the remote file is not a directory \" "+ Remote +" "\" skip download ")
		return c.downloadDir(remote, local)

	}
	return c.downloadFile(remote, local)

}

//downloadFile a file from the remote server like cp
func (c *Client) downloadFile(remoteFile, local string) error {
	//remoteFile = strings.TrimSuffix(remoteFile, "/")
	if !c.IsFile(remoteFile) {
		return errors.New("sftp:  The file does not exist or is not a file ,  Skip directory Download  downloadFile(" + remoteFile + ")")
	}
	var localFile string
	if local[len(local)-1] == '/' {
		localFile = filepath.Join(local, filepath.Base(remoteFile))
	} else {
		localFile = local
	}
	localFile = filepath.ToSlash(localFile)
	if c.Size(remoteFile) > 1000 {
		rsum := c.Md5File(remoteFile)
		ioutil.WriteFile(localFile+".md5", []byte(rsum), 755)
		if FileExist(localFile) {
			if rsum != "" {
				lsum, _ := Md5File(localFile)
				if lsum == rsum {
					log.Println("sftp:  File is consistent with local ， Skip Download ！", localFile)
					return nil
				}
				log.Println("sftp:  Downloading  ", localFile)
			}
		}
	}
	if err := os.MkdirAll(filepath.Dir(localFile), os.ModePerm); err != nil {
		//log.Println(err)
		return err
	}

	r, err := c.SFTPClient.Open(remoteFile)
	if err != nil {
		return err
	}
	defer r.Close()

	l, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer l.Close()

	_, err = io.Copy(l, r)
	return err
}

//downloadDir from remote dir to local dir like rsync
//rsync -av src/ dst     ./src/* --> /root/dst/*
//rsync -av src/ dst/    ./src/* --> /root/dst/*
//rsync -av src  dst     ./src/* --> /root/dst/src/*
//rsync -av src  dst/    ./src/* --> /root/dst/src/*
func (c *Client) downloadDir(remote, local string) error {
	var localDir, remoteDir string

	if !c.IsDir(remote) {
		return errors.New("sftp:  The directory does not exist or is not a directory ,  skip  downloadDir(" + remote + ")")
	}
	remoteDir = remote
	if remote[len(remote)-1] == '/' {
		localDir = local
	} else {
		localDir = path.Join(local, path.Base(remote))
	}

	walker := c.SFTPClient.Walk(remoteDir)

	for walker.Step() {
		if err := walker.Err(); err != nil {
			log.Println(err)
			continue
		}

		info := walker.Stat()

		relPath, err := filepath.Rel(remoteDir, walker.Path())
		if err != nil {
			return err
		}

		localPath := filepath.ToSlash(filepath.Join(localDir, relPath))

		//if we have something at the download path delete it if it is a directory
		//and the remote is a file and vice a versa
		localInfo, err := os.Stat(localPath)
		if os.IsExist(err) {
			if localInfo.IsDir() {
				if info.IsDir() {
					continue
				}

				err = os.RemoveAll(localPath)
				if err != nil {
					return err
				}
			} else if info.IsDir() {
				err = os.Remove(localPath)
				if err != nil {
					return err
				}
			}
		}

		if info.IsDir() {
			err = os.MkdirAll(localPath, os.ModePerm)
			if err != nil {
				return err
			}

			continue
		}

		c.downloadFile(walker.Path(), localPath)

	}
	return nil
}

//UploadFile upload local file to SFTP remote directory remote
func (c *Client) UploadFile(localFile, remote string) error {
	//localFile = strings.TrimSuffix(localFile, "/")
	//localFile = filepath.ToSlash(localFile)
	info, err := os.Stat(localFile)
	if err != nil || info.IsDir() {
		return errors.New("sftp:  The local file does not exist , Or is it a file  UploadFile(\"" + localFile + "\")  Skip upload ")
	}

	l, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer l.Close()

	var remoteFile, remoteDir string
	if remote[len(remote)-1] == '/' {
		remoteFile = filepath.ToSlash(filepath.Join(remote, filepath.Base(localFile)))
		remoteDir = remote
	} else {
		remoteFile = remote
		remoteDir = filepath.ToSlash(filepath.Dir(remoteFile))
	}
	log.Println("sftp: UploadFile", localFile, remoteFile)
	if info.Size() > 1000 {
		//1. Check whether the remote exists
		rsum := c.Md5File(remoteFile)
		if rsum != "" {
			lsum, _ := Md5File(localFile)
			if lsum == rsum {
				log.Println("sftp:  File is consistent with local ， Skip upload ！", localFile)
				return nil
			}
			log.Println("sftp:  Uploading  ", localFile)
		}
	}

	//If the directory does not exist, remotedir is created
	if _, err := c.SFTPClient.Stat(remoteDir); err != nil {
		log.Println("sftp: Mkdir all", remoteDir)
		c.MkdirAll(remoteDir)
	}

	r, err := c.SFTPClient.Create(remoteFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(r, l)
	return err
}

//UploadDir files without checking diff status
func (c *Client) UploadDir(localDir string, remoteDir string) (err error) {
	//defer func() {
	//if err !=  nil {
	//err = errors.New("UploadDir " + err.Error())
	//}
	//}()
	//Local input detection, must be directory
	//localDir = filepath.ToSlash(localDir)
	info, err := os.Stat(localDir)
	if err != nil || !info.IsDir() {
		return errors.New("sftp:  The local directory does not exist or is not a directory  UploadDir(\"" + localDir + "\")  Skip upload ")
	}

	//If Rsync localdir does not end with '/', a tail directory is created
	if localDir[len(localDir)-1] != '/' {
		remoteDir = filepath.ToSlash(filepath.Join(remoteDir, filepath.Base(localDir)))
	}
	log.Println("sftp: UploadDir", localDir, remoteDir)

	rootDst := strings.TrimSuffix(remoteDir, "/")
	if c.IsFile(rootDst) {
		c.SFTPClient.Remove(rootDst)
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//Calculate the final destination using the
		//base source and root destination
		relSrc, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}
		finalDst := filepath.Join(rootDst, relSrc)

		//In Windows, Join uses backslashes which we don't want to get
		//to the sftp server
		finalDst = filepath.ToSlash(finalDst)

		//Skip the creation of the target destination directory since
		//it should exist and we might not even own it
		if finalDst == remoteDir {
			return nil
			log.Println("sftp: ", remoteDir, "--->", finalDst)

		}

		if info.IsDir() {
			err := c.MkdirAll(finalDst)
			if err != nil {
				log.Println("sftp: MkdirAll", err)
			}
			//log.Println("MkdirAll", finalDst)
			//err = c.SFTPClient.Mkdir(finalDst)
			//log.Println(err)
			//if err := c.SFTPClient.Mkdir(finalDst);  err !=  nil {
			// 	// Do not consider it an error if the directory existed
			//remoteFi, fiErr := c.SFTPClient.Lstat(finalDst)
			//if fiErr !=  nil || ! remoteFi.IsDir() {
			//return err
			//}
			//}
			//return err
		} else {
			//f, err := os.Open(path)
			//if err !=  nil {
			//return err
			//}
			//defer f.Close()
			return c.UploadFile(path, finalDst)
		}
		return nil

	}
	return filepath.Walk(localDir, walkFunc)
}

//Remove a file from the remote server
func (c *Client) Remove(path string) error {
	return c.SFTPClient.Remove(path)
}

//RemoveDirectory Remove a directory from the remote server
func (c *Client) RemoveDirectory(path string) error {
	return c.SFTPClient.RemoveDirectory(path)
}

//ReadAll Read a remote file and return the contents.
func (c *Client) ReadAll(filepath string) ([]byte, error) {
	file, err := c.SFTPClient.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

//Does the fileexist file exist
func (c *Client) FileExist(filepath string) (bool, error) {
	if _, err := c.SFTPClient.Stat(filepath); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) RemoveFile(remoteFile string) error {
	return c.SFTPClient.Remove(remoteFile)
}
func (c *Client) RemoveDir(remoteDir string) error {
	remoteFiles, err := c.SFTPClient.ReadDir(remoteDir)
	if err != nil {
		log.Printf("sftp: remove remote dir: %s err: %v\n", remoteDir, err)
		return err
	}
	for _, file := range remoteFiles {
		subRemovePath := path.Join(remoteDir, file.Name())
		if file.IsDir() {
			c.RemoveDir(subRemovePath)
		} else {
			c.RemoveFile(subRemovePath)
		}
	}
	c.SFTPClient.RemoveDirectory(remoteDir) //must empty dir to remove
	log.Printf("sftp: remove remote dir: %s ok\n", remoteDir)
	return nil
}

//RemoveAll recursively deletes directories and files
func (c *Client) RemoveAll(remoteDir string) error {
	c.RemoveDir(remoteDir)
	return nil
}

//Mkdirall creates directories recursively
func (c *Client) MkdirAll(dirpath string) error {

	parentDir := filepath.ToSlash(filepath.Dir(dirpath))
	_, err := c.SFTPClient.Stat(parentDir)
	if err != nil {
		//log.Println(err)
		if err.Error() == "file does not exist" {
			err := c.MkdirAll(parentDir)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	err = c.SFTPClient.Mkdir(filepath.ToSlash(dirpath))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Mkdir(path string, fi os.FileInfo) error {
	log.Printf("[DEBUG] sftp: creating dir %s", path)

	if err := c.SFTPClient.Mkdir(path); err != nil {
		//Do not consider it an error if the directory existed
		remoteFi, fiErr := c.SFTPClient.Lstat(path)
		if fiErr != nil || !remoteFi.IsDir() {
			return err
		}
	}

	mode := fi.Mode().Perm()
	if err := c.SFTPClient.Chmod(path, mode); err != nil {
		return err
	}
	return nil
}

//Isdir checks if the remote is a directory
func (c *Client) IsDir(path string) bool {
	//Check whether the remote is a file or directory
	info, err := c.SFTPClient.Stat(path)
	if err == nil && info.IsDir() {
		return true
	}
	return false
}

//Size gets the file size
func (c *Client) Size(path string) int64 {
	info, err := c.SFTPClient.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

//Isfile checks if the remote is a file
func (c *Client) IsFile(path string) bool {
	info, err := c.SFTPClient.Stat(path)
	if err == nil && !info.IsDir() {
		return true
	}
	return false
}

//Isnotexist checks if the remote file does not exist
func (c *Client) IsNotExist(path string) bool {
	_, err := c.SFTPClient.Stat(path)
	return err != nil
}

//Isexist checks if the remote file exists
func (c *Client) IsExist(path string) bool {

	_, err := c.SFTPClient.Stat(path)
	return err == nil
}

//Md5file check if the remote file exists
func (c *Client) Md5File(path string) string {
	if c.IsNotExist(path) {
		return ""
	}
	b, err := c.Output("md5sum " + path)
	if err != nil {
		return ""
	}
	return string(bytes.Split(b, []byte{' '})[0])

}
