GOPATH=$(shell go env GOPATH)
PRODUCT=mealie-crypt
FULL_PRODUCT=github.com/warrenhodg/${PRODUCT}
GOLANG_DOCKER_IMAGE="golang:1.11"
GO_SRC=${GOPATH}/src

${PRODUCT}: windows linux mac
windows: ${PRODUCT}.exe
linux: ${PRODUCT}.linux
mac: ${PRODUCT}.mac

help:
	@echo "Build normally"
	@echo "\tmake"
	@echo "Build using docker"
	@echo "\tUSE_DOCKER=1 make"

${PRODUCT}.linux: *.go
ifeq (${USE_DOCKER}, 1)
	docker run --rm -v ${PWD}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} ${GOLANG_DOCKER_IMAGE} make linux
else
	GOOS=linux GOARCH=386 go build -o ${PRODUCT}.linux
endif

${PRODUCT}.mac: *.go
ifeq (${USE_DOCKER}, 1)
	docker run --rm -v ${PWD}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} ${GOLANG_DOCKER_IMAGE} make mac
else
	GOOS=darwin GOARCH=amd64 go build -o ${PRODUCT}.mac
endif

${PRODUCT}.exe: *.go
ifeq (${USE_DOCKER}, 1)
	docker run --rm -v ${PWD}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} ${GOLANG_DOCKER_IMAGE} make windows
else
	GOOS=windows GOARCH=386 go build -o ${PRODUCT}.exe
endif

clean:
	rm -f ${PRODUCT}.exe ${PRODUCT}.linux ${PRODUCT}.mac
