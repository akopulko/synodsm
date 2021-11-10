BINARY_NAME=synodsm

build:
	GOARCH=arm64 GOOS=darwin go build -o ${BINARY_NAME} app/*.go

clean:
	go clean
	rm ${BINARY_NAME}