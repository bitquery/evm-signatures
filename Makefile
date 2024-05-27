# Code test
.PHONY: test
test:
	go test -race -count=1 -shuffle=on -failfast -timeout=30s ./...

# Code tidy
.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...