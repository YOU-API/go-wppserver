BINARY_NAME=wppserver

compile:
	echo "Compiling for every OS and Platform"
	#GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-windows-386.exe cmd/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-windows-amd64.exe cmd/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-darwin-amd64 cmd/main.go
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-386 cmd/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-amd64 cmd/main.go
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-arm cmd/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-arm64 cmd/main.go
	#GOOS=linux GOARCH=riscv go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-riscv cmd/main.go
	GOOS=linux GOARCH=riscv64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-riscv64 cmd/main.go
	GOOS=freebsd GOARCH=arm64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-arm64 cmd/main.go
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-amd64 cmd/main.go
	#GOOS=freebsd GOARCH=riscv64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-riscv64 cmd/main.go
	GOOS=openbsd GOARCH=arm64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-arm64 cmd/main.go
	GOOS=openbsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-amd64 cmd/main.go
	#GOOS=openbsd GOARCH=riscv64 go build -ldflags="-s -w" -o bin/${BINARY_NAME}-linux-riscv64 cmd/main.go