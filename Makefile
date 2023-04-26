default: build

build: out/craftctl

out/craftctl:
	go build -ldflags="-X 'main.Version=1.0-SNAPSHOT'" -o out/craftctl ./cmd/

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf out
