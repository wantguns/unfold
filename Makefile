run:
	go run main.go

build:
	go build -o bin/unfold main.go

check-build:
	# Invoke goreleaser
	
	goreleaser check
	goreleaser release --snapshot --clean
