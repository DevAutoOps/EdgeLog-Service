package notice

import "os"

type INotice interface {
	SendText(to []string, str string) error
	SendFile(to []string, file *os.File) error
}
