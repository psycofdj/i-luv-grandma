all: test

test:
	@go test -v ./...

check:
	@golangci-lint run --config .golangci.yml

coverage:
	@go test -cover -coverprofile cover.out -v ./...
	@go tool cover -func=cover.out
	@rm -f cover.out

perfs:
	@go build -o i-luv-grandma .
	@./i-luv-grandma -profile output.pprof -input dataset/4320p.pbm -output /dev/null -angle 180
	@go tool pprof -top i-luv-grandma output.pprof
