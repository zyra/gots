APP_VERSION := $(shell git tag | tail -1)

.PHONY: build
build: build_linux build_windows build_darwin ; @echo "Done building!"

build_linux: ; @\
GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "-X main.AppVersion=${APP_VERSION}" -o bin/gots_linux_amd64 cmd/gots/main.go && \
chmod +x bin/gots_linux_amd64

build_windows: ; @\
GOOS=windows GOARCH=amd64 go build -mod vendor -ldflags "-X main.AppVersion=${APP_VERSION}" -o bin/gots_windows_amd64.exe cmd/gots/main.go

build_darwin: ; @\
GOOS=darwin GOARCH=amd64 go build -mod vendor -ldflags "-X main.AppVersion=${APP_VERSION}" -o bin/gots_darwin_amd64 cmd/gots/main.go && \
chmod +x bin/gots_darwin_amd64

docker_build: ; @\
docker build -t harbor.zyra.ca/public/gots .

docker_push: ; @\
docker push harbor.zyra.ca/public/gots