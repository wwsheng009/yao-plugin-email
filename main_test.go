package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/yaoapp/kun/grpc"
)

func TestEmailPlugin_Exec(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %s", err.Error())
	}
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailUser := os.Getenv("EMAIL_USERNAME")

	requestBody := fmt.Sprintf(`{
		"account":{
			"server":"smtp.qq.com", 
			"port":587,
			"username": "%s", 
			"password": "%s",
			"type":"stmp"
		},
		"messages":[
			{
				"from": "%s",
				"to":[{"Address":"%s"}],
				"cc":[{"Address":"%s"}],
				"subject": "小佩奇",
				"body": "<h1>新年快乐</h1>",
				"attachments": ["./yaoapp/data/test.jpg"]
			}
		]
	}`, emailUser, emailPassword, emailUser, emailUser, emailUser)
	type args struct {
		name string
		args []interface{}
	}
	tests := []struct {
		name    string
		plugin  *EmailPlugin
		args    args
		want    *grpc.Response
		wantErr bool
	}{
		{
			name:   "test",
			plugin: &EmailPlugin{},
			args: struct {
				name string
				args []interface{}
			}{
				name: "send",
				args: []interface{}{requestBody},
			},
			want: &grpc.Response{Bytes: []byte(`{"code":200,"message":"操作成功"}`), Type: "map"},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.Exec(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailPlugin.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailPlugin.Exec() = %v, want %v", string(got.Bytes), string(tt.want.Bytes))
			}
		})
	}
}

func TestEmailPlugin_Exec2(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %s", err.Error())
	}
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailUser := os.Getenv("EMAIL_USERNAME")

	requestBody := fmt.Sprintf(`{
		"account":{
			"server":"smtp.gmail.com", 
			"port":587,
			"username": "%s", 
			"password": "%s",
			"type":"smtp"
		},
		"messages":[
			{
				"from": "%s",
				"to":[{"Address":"%s"}],
				"cc":[{"Address":"%s"}],
				"subject": "小佩奇",
				"body": "<h1>新年快乐</h1>",
				"attachments": ["./yaoapp/data/test.jpg"]
			}
		]
	}`, emailUser, emailPassword, emailUser, emailUser, emailUser)
	type args struct {
		name string
		args []interface{}
	}
	tests := []struct {
		name    string
		plugin  *EmailPlugin
		args    args
		want    *grpc.Response
		wantErr bool
	}{
		{
			name:   "test gmail smtp account",
			plugin: &EmailPlugin{},
			args: struct {
				name string
				args []interface{}
			}{
				name: "send",
				args: []interface{}{requestBody},
			},
			want: &grpc.Response{Bytes: []byte(`{"code":200,"message":"发送成功"}`), Type: "map"},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.Exec(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailPlugin.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailPlugin.Exec() = %v, want %v", string(got.Bytes), string(tt.want.Bytes))
			}
		})
	}
}

func TestEmailPluginReceive_Exec(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %s", err.Error())
	}
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailUser := os.Getenv("EMAIL_USERNAME")

	requestBody := fmt.Sprintf(`{
		"account":{
			"server":"imap.qq.com", 
			"port":993,
			"username": "%s", 
			"password": "%s",
			"type":"imap"
		}
	}`, emailUser, emailPassword)
	type args struct {
		name string
		args []interface{}
	}
	tests := []struct {
		name    string
		plugin  *EmailPlugin
		args    args
		want    *grpc.Response
		wantErr bool
	}{
		{
			name:   "test",
			plugin: &EmailPlugin{},
			args: struct {
				name string
				args []interface{}
			}{
				name: "receive",
				args: []interface{}{requestBody},
			},
			want: &grpc.Response{Bytes: []byte(`{"code":200,"message":"接收成功"}`), Type: "map"},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.Exec(tt.args.name, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("EmailPlugin.Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailPlugin.Exec() = %v, want %v", string(got.Bytes), string(tt.want.Bytes))
			}
		})
	}
}
