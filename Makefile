PRODUCT=mealie-crypt
FULL_PRODUCT=github.com/warrenhodg/${PRODUCT}
GOLANG_DOCKER_IMAGE="golang:1.10"
GO_SRC=${GOPATH}/src
DEP=${GOPATH}/bin/dep
VENDOR=vendor

${PRODUCT}: windows linux mac
windows: ${PRODUCT}.exe
linux: ${PRODUCT}.linux
mac: ${PRODUCT}.mac

help:
	@echo "Build normally"
	@echo "\tmake"
	@echo "Build using docker"
	@echo "\tUSE_DOCKER=1 make"

${PRODUCT}.linux: *.go ${VENDOR}
ifeq (${USE_DOCKER}, 1)
	docker run --rm -v ${PWD}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} golang:1.10 make linux
else
	GOOS=linux GOARCH=386 go build -o ${PRODUCT}.linux ${FULL_PRODUCT}
endif

${PRODUCT}.mac: *.go ${VENDOR}
ifeq (${USE_DOCKER}, 1)
	docker run --rm -v ${PWD}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} golang:1.10 make mac
else
	GOOS=darwin GOARCH=386 go build -o ${PRODUCT}.mac ${FULL_PRODUCT}
endif

${PRODUCT}.exe: *.go ${VENDOR}
ifeq (${USE_DOCKER}, 1)
	docker run --rm -v ${PWD}:${GO_SRC}/${FULL_PRODUCT} -w ${GO_SRC}/${FULL_PRODUCT} golang:1.10 make windows
else
	GOOS=windows GOARCH=386 go build -o ${PRODUCT}.exe ${FULL_PRODUCT}
endif

${VENDOR}: ${DEP}
ifneq (${USE_DOCKER}, 1)
	dep ensure
	chmod 777 vendor
endif

${DEP}:
ifneq (${USE_DOCKER}, 1)
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

clean:
	rm -f ${PRODUCT}.exe ${PRODUCT}.linux ${PRODUCT}.mac
