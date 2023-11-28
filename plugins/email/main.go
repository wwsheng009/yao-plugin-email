package main

//插件模板
import (
	"crypto/tls"
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/yaoapp/kun/grpc"
	"gopkg.in/gomail.v2"
)

// 定义插件类型，包含grpc.Plugin
type EmailPlugin struct{ grpc.Plugin }

type Account struct {
	Server   string `json:"server"`
	Port     int    `json:"port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type Message struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	CC          []string `json:"cc,omitempty"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Attachments []string `json:"attachments,omitempty"`
}

type Email struct {
	Account  Account   `json:"account"`
	Messages []Message `json:"messages"`
}

// 设置插件日志到单独的文件
func (plugin *EmailPlugin) setLogFile() {
	var output io.Writer = os.Stdout
	//开启日志
	logroot := os.Getenv("GOU_TEST_PLG_LOG")
	if logroot == "" {
		logroot = "./logs"
	}
	if logroot != "" {
		logfile, err := os.OpenFile(path.Join(logroot, "email.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			output = logfile
		}
	}
	plugin.Plugin.SetLogger(output, grpc.Trace)
}

func sendEmails(email Email) []error {
	errs := make([]error, 0)
	for _, v := range email.Messages {
		err := sendEmail(email.Account, v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
func sendEmail(account Account, message Message) error {
	m := gomail.NewMessage()
	m.SetHeader("From", message.From)
	if message.From == "" {
		m.SetHeader("From", account.Username)
	}
	m.SetHeader("To", message.To)
	if len(message.CC) > 0 {
		m.SetHeader("Cc", message.CC...)
	}
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/html", message.Body)

	for _, attachment := range message.Attachments {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(account.Server, account.Port, account.Username, account.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// 插件执行需要实现的方法
// 参数name是在调用插件时的方法名，比如调用插件demo的Hello方法是的规则是plugins.demo.Hello时。
//
// 注意：name会自动的变成小写
//
// args参数是一个数组，需要在插件中自行解析。判断它的长度与类型，再转入具体的go类型。
//
// Exec 插件入口函数

func (plugin *EmailPlugin) Exec(name string, args ...interface{}) (*grpc.Response, error) {
	// plugin.Logger.Log(hclog.Trace, "plugin method called", name)
	// plugin.Logger.Log(hclog.Trace, "args", args)
	isOk := true
	var v = make(map[string]interface{})
	var email Email

	switch name {
	case "send":
		if len(args) < 1 {
			v = map[string]interface{}{"code": 400, "message": "参数不足，需要一个参数"}
			isOk = false
		}
		if isOk {
			switch data := args[0].(type) {

			case string:
				err := json.Unmarshal([]byte(data), &email)
				if err != nil {
					isOk = false
					v = map[string]interface{}{"code": 400, "message": err.Error()}
				}
			case map[string]interface{}:
				jsonData, err := json.Marshal(data)
				if err != nil {
					isOk = false
					v = map[string]interface{}{"code": 400, "message": err.Error()}
				}
				err = json.Unmarshal(jsonData, &email)
				if err != nil {
					isOk = false
					v = map[string]interface{}{"code": 400, "message": err.Error()}
				}
			default:
				isOk = false
				v = map[string]interface{}{"code": 400, "message": "传入参数类型错误,请传入json数据"}
			}

			if isOk {
				errs := sendEmails(email)
				if len(errs) > 0 {
					isOk = false
					message := ""
					for _, e := range errs {
						message += e.Error() + "\n"
					}
					v = map[string]interface{}{"code": 400, "message": message}
				}
			}

		}

	}
	if isOk {
		v = map[string]interface{}{"code": 200, "message": "发送成功"}
	}

	//输出前需要转换成字节
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	//设置输出数据的类型
	//支持的类型：map/interface/string/integer,int/float,double/array,slice
	return &grpc.Response{Bytes: bytes, Type: "map"}, nil
}

// 生成插件时函数名修改成main
func main() {
	plugin := &EmailPlugin{}
	plugin.setLogFile()
	grpc.Serve(plugin)
}
