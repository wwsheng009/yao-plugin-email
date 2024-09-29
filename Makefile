
build:
	CGO_ENABLED=0 go build -o yaoapp/plugins/email.so
	chmod +x yaoapp/plugins/email.so
	
windows:
	set CGO_ENABLED=0
	go build -o yaoapp/plugins/email.dll

.PHONY:	clean
clean:
	rm -f yaoapp/plugins/email.so