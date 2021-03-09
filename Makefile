build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rd-tool ./main.go
build-macos-x64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o rd-tool-mac ./main.go
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o rd-tool.exe ./main.go