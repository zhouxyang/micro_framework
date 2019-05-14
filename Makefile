MAIN_VER=$(shell awk -F '"' '$$1=="MainVersion = " {print $$2;exit}' changelog.md)
GIT_CNT=$(shell git rev-list --count HEAD)
VERSION=${MAIN_VER}.${GIT_CNT}
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
DATE=$(shell date "+%Y-%m-%d %H:%M:%S")

all:  clean server client

server:
	go build -i -v -ldflags "-X 'main.version=version: ${VERSION}, date: ${DATE}'" -o server/server-${BRANCH}-v${VERSION} server/server.go

client:
	go build -i -v -ldflags "-X 'main.version=version: ${VERSION}, date: ${DATE}'" -o client/client-${BRANCH}-v${VERSION} client/client.go


clean:
	rm -f server/server-${BRANCH}-v${VERSION}
	rm -f client/client-${BRANCH}-v${VERSION}

.PHONY: all server client clean


