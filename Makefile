
.PHONY: build
build: clean
	@go build -o mdtoc cmd/mdtoc/main.go

.PHONY: clean
clean:
	@rm -f mdtoc

.PHONY: test
test:
	@go test -v ./...
