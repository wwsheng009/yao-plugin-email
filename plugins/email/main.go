package main

//插件模板
import (
	"encoding/json"
	"io"
	"os"
	"path"

	"github.com/yaoapp/kun/grpc"
)

// 定义插件类型，包含grpc.Plugin
type EmailPlugin struct{ grpc.Plugin }

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
	var v = map[string]interface{}{"code": 200, "message": "发送成功"}
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
				// if email.Account.Type != "imap" {
				errs := sendStmpEmails(email)
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
	case "receive":
		if len(args) < 1 {
			v = map[string]interface{}{"code": 400, "message": "参数不足，需要一个参数"}
			isOk = false
		}

		if isOk {

			// var account Account
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
				messages, err := receiveImapEmails(email)
				if err != nil {
					isOk = false
					v = map[string]interface{}{"code": 400, "message": err.Error()}
				} else {
					v = map[string]interface{}{"code": 200, "emails": messages}
				}
			}

		}
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
