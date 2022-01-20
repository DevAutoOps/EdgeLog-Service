package tools

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

//  Calculation file MD5 value
func FileMD5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", err)
		return ""
	}

	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ReadAll", err)
		return ""
	}

	return fmt.Sprintf("%x", md5.Sum(body))
}

//  Calculation file MD5 value
func FileMD5v2(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open", err)
		return ""
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return ""
	}

	return fmt.Sprintf("%x", md5hash.Sum(nil))
}

//  Calculate the content of the file MD5 value
func FileMD5v3(content string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(content)))
}

//  Delete directory
func RemoveDir(path string) {
	//  Delete file
	dir := path
	exist, err := PathExists(dir)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if exist {
			// os.RemoveAll  Is traversal delete ， Both folders and files can be used
			err := os.RemoveAll(dir)
			if err != nil {
				fmt.Println(dir+" Deletion failed ：", err.Error())
			} else {
				fmt.Println(dir + " Delete succeeded ！")
			}
		} else {
			fmt.Println(dir + " file 、 Folder does not exist ！")
		}
	}
}

//  write file
func WriteFile(path string, content string) error {
	err := ioutil.WriteFile(path, []byte(content), 0777)
	if err != nil {
		return err
	}
	return nil
}
