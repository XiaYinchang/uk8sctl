all: build

repo=main
commit=`git rev-parse HEAD`
exe=uk8sctl

.PHONY: build
build:
	@echo $(object)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -mod=vendor -o ./$(exe) $(repo).go