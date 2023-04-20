default: build

build: out/mccp

out/mccp:
	go build -ldflags="-X 'main.Version=1.0-SNAPSHOT'" -o out/mccp ./cmd/

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf out
