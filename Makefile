BIN=gc-sso
BIN_LINUX=$(BIN)-linux

.PHONY: local
local:
	go build -o bin/$(BIN) .

.PHONY: linux
linux:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w"  -o bin/$(BIN_LINUX) .

.PHONY: image
image:
	docker build -t mivinci/$(BIN) .

.PHONY: run
run:
	docker run --rm -p 8000:8000 mivinci/$(BIN)

.PHONY: rm
rm:
	rm -r data/*