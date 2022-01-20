package ssh

import (
	"testing"
)

func TestClient_IsCheck(t *testing.T) {
	c := GetClient()
	defer c.Close()
	var remotes = []string{
		"/root/test/notExist",
		"/root/test/notExist/",
		"/root/test/file",
		"/root/test/file/", //non-existent
		"/root/test/dir",
		"/root/test/dir/",
	}

	///root/test/file  		 existence
	///root/test/file/  	 non-existent
	///root/test/dir  		 existence
	///root/test/dir/  		 existence
	for _, v := range remotes {
		is := c.IsExist(v)
		if is {
			println(v, "\t existence ")
		} else {
			println(v, "\t non-existent ")
		}
	}

	///root/test/file  		 Not a directory
	///root/test/file/  	 Not a directory
	///root/test/dir  		 Is a directory
	///root/test/dir/  		 Is a directory
	println()
	for _, v := range remotes {
		isdir := c.IsDir(v)
		if isdir {
			println(v, "\t Is a directory ")
		} else {
			println(v, "\t Not a directory ")
		}
	}

	///root/test/file  		 It's a file
	///root/test/file/  	 Not a file
	///root/test/dir  		 Not a file
	///root/test/dir/  		 Not a file
	println()
	for _, v := range remotes {
		isfile := c.IsFile(v)
		if isfile {
			println(v, "\t It's a file ")
		} else {
			println(v, "\t Not a file ")
		}

	}

}
