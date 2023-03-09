OS = linux
ARCH = amd64

BUILD_FLAGS = -x -v
BUILD_OUTPUT = unitedb-discord-bot

.PHONY: run all tidy

all: build

tidy:
	go mod tidy

run: tidy
	go run ./cmd/bot

build: tidy
	env GOOS=${OS} GOARCH=${ARCH} go build ${BUILD_FLAGS} -o ${BUILD_OUTPUT} ./cmd/bot
