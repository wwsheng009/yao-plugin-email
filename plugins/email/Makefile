
build:
	CGO_ENABLED=0 go build -o ../email.so

windows:
	set CGO_ENABLED=0
	go build -o ../email.dll

.PHONY:	clean
clean:
	rm -f "../email.so"