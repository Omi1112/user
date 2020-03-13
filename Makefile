TEST_FILES=$(shell find -name '*_test.go')
build: task
	go version
task:
	go test ./...
