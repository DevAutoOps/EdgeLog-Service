package notice

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	LF = "\r\n"
)

type EmailNotice struct {
	recipient []string
	user      string
	subject   string
	client    *smtp.Client
}

func CreateEmailNotice(s SMTP) (INotice, error) {
	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := tls.DialWithDialer(&dialer, "tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port), &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.Addr,
	})
	// conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.Addr, s.Port), &tls.Config{
	// 	InsecureSkipVerify: true,
	// 	ServerName:         s.Addr,
	// })
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, s.Addr)
	if err != nil {
		return nil, err
	}
	auth := smtp.PlainAuth("", s.User, s.Pass, s.Addr)
	if ok, _ := client.Extension("AUTH"); ok {
		if err = client.Auth(auth); err != nil {
			return nil, err
		}
	}
	split := strings.Split(s.ReceiveEmail, ",")
	return &EmailNotice{
		client:    client,
		subject:   s.Topic,
		user:      s.User,
		recipient: split,
	}, nil
}

func (s *EmailNotice) makeMsg(subject string, to []string, value string,
	withFile bool, fileContentType, fileName string, file []byte) []byte {
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	headerMap := make(map[string]string)
	headerMap["From"] = s.user
	headerMap["To"] = strings.Join(to, ";")
	//headerMap["Cc"] = strings.Join(cc, ";")
	//headerMap["Bcc"] = strings.Join(bcc, ";")
	headerMap["Subject"] = subject
	headerMap["Content-Type"] = "multipart/mixed;boundary=" + boundary
	headerMap["Mime-Version"] = "1.0"
	headerMap["Date"] = time.Now().String()
	headerStr := ""
	for key, value := range headerMap {
		headerStr += key + ":" + value + LF
	}
	headerStr += LF
	buffer.WriteString(headerStr)
	contentType := "text/plain;charset=utf-8"
	body := LF + "--" + boundary + LF
	body += "Content-Type:" + contentType + LF
	body += LF + value + LF
	buffer.WriteString(body)
	if withFile {
		attachment := LF + "--" + boundary + LF
		attachment += "Content-Transfer-Encoding:base64" + LF
		attachment += "Content-Disposition:attachment" + LF
		attachment += "Content-Type:" + fileContentType + ";name=\"" + fileName + "\"" + LF
		buffer.WriteString(attachment)
		payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
		base64.StdEncoding.Encode(payload, file)
		buffer.WriteString(LF)
		for index, line := 0, len(payload); index < line; index++ {
			buffer.WriteByte(payload[index])
			if (index+1)%76 == 0 {
				buffer.WriteString(LF)
			}
		}
	}
	buffer.WriteString(LF + "--" + boundary + "--")
	return buffer.Bytes()
}

func (s *EmailNotice) SendText(to []string, msg string) (err error) {
	to = append(to, s.recipient...)
	content := s.makeMsg(s.subject, to, msg,
		false, "", "", []byte{})
	return s.Send(to, content)
}

func (s *EmailNotice) SendFile(to []string, file *os.File) (err error) {
	fileValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	to = append(to, s.recipient...)
	content := s.makeMsg(s.subject, to, "",
		true, "application/octet-stream", file.Name(), fileValue)
	return s.Send(to, content)
}

func (s *EmailNotice) Send(to []string, content []byte) (err error) {
	if err = s.client.Mail(s.user); err != nil {
		return fmt.Errorf("Client:clientMail:%v", err)
	}
	for _, addr := range to {
		if err = s.client.Rcpt(addr); err != nil {
			return fmt.Errorf("Client:Rcpt:%v", err)
		}
	}
	w, err := s.client.Data()
	if err != nil {
		return fmt.Errorf("Client:Data:%v", err)
	}
	_, err = w.Write(content)
	if err != nil {
		return fmt.Errorf("Client:WriterBody:%v", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("Client:CloseBody:%v", err)
	}
	return
}
