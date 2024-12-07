.PHONY: default
default: build
build:
	CGO_ENABLED=0 go build -o bin/wg-trick-server server.go

install:
	cp bin/wg-trick-server /usr/local/bin/ && \
	[ -d /etc/systemd/system/ ] && cp wg-trick-server.service /etc/systemd/system/