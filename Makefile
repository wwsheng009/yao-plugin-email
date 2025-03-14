
build:
	CGO_ENABLED=0 go build -o yaoapp/plugins/email.so
	chmod +x yaoapp/plugins/email.so
	
windows:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o yaoapp/plugins/email.dll

.PHONY:	clean
clean:
	rm -f yaoapp/plugins/email.*