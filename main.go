package main

//插件模板
import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

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

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 创建错误响应
func newErrorResponse(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// 创建成功响应
func newSuccessResponse(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "操作成功",
		Data:    data,
	}
}

// 解析参数
func parseArgs(args []interface{}) (*Email, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("参数不足，需要一个参数")
	}

	var email Email
	switch data := args[0].(type) {
	case string:
		if err := json.Unmarshal([]byte(data), &email); err != nil {
			return nil, err
		}
	case map[string]interface{}:
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(jsonData, &email); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("传入参数类型错误,请传入json数据")
	}

	return &email, nil
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
	var response Response

	switch name {
	case "send":
		email, err := parseArgs(args)
		if err != nil {
			response = newErrorResponse(400, err.Error())
		} else {
			if plugin.Plugin.Logger.IsDebug() {
				// 输出转换后的email结构
				jsonBytes, err := json.MarshalIndent(email, "", "    ")
				if err != nil {
					return nil, fmt.Errorf("序列化email失败: %v", err)
				}
				plugin.Plugin.Logger.Info("转换后的email结构: %s", string(jsonBytes))
			}

			errs := sendStmpEmails(*email)
			if len(errs) > 0 {
				var message strings.Builder
				for _, e := range errs {
					message.WriteString(e.Error())
					message.WriteString("\n")
				}
				response = newErrorResponse(400, message.String())
			} else {
				response = newSuccessResponse(nil)
			}
		}

	case "receive":
		email, err := parseArgs(args)
		if err != nil {
			response = newErrorResponse(400, err.Error())
		} else {
			if plugin.Plugin.Logger.IsDebug() {
				// 输出转换后的email结构
				jsonBytes, err := json.MarshalIndent(email, "", "    ")
				if err != nil {
					return nil, fmt.Errorf("序列化email失败: %v", err)
				}
				plugin.Plugin.Logger.Info("转换后的email结构: %s", string(jsonBytes))
			}
			messages, err := receiveImapEmails(*email)
			if err != nil {
				response = newErrorResponse(400, err.Error())
			} else {
				response = newSuccessResponse(map[string]interface{}{"emails": messages})
			}
		}

	default:
		response = newErrorResponse(404, fmt.Sprintf("未找到方法: %s", name))
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return &grpc.Response{Bytes: bytes, Type: "map"}, nil
}

// 生成插件时函数名修改成main
func main() {
	plugin := &EmailPlugin{}
	plugin.setLogFile()
	grpc.Serve(plugin)
}
