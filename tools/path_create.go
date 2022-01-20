package tools

import (
	"errors"
	"os"
)

var AlreadyExists = " The folder or file already exists "

//  Determine whether the folder exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//  create folder
func MkdirFilePath(path string) (string, error) {
	err := os.Mkdir(path, 0777)
	if err != nil {
		return "", err
	}
	return "", nil
}

func IsExists(path string, isFile bool) bool {
	s, err := os.Stat(path)
	if err == nil {
		if s.IsDir() == !isFile {
			return true
		}
	}
	return false
}

func CreateFile(filePath string, content []byte) error {
	if IsExists(filePath, true) {
		return errors.New(AlreadyExists)
	}
	newFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = newFile.Write(content)
	return err
}

func CreateFolder(filePath string) error {
	if IsExists(filePath, false) {
		return errors.New(AlreadyExists)
	}
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		return err
	}
	err = os.Chmod(filePath, 0777)
	if err != nil {
		return err
	}
	return nil
}
