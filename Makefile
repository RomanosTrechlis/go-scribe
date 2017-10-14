PROJECT = logScribe
GITHUB = /home/romanos/go/src/github.com/RomanosTrechlis
PROJECT_DIR = ${GITHUB}/${PROJECT}
LOG_SCRIBE_CMD = ${PROJECT_DIR}/cmd/${PROJECT}
MEDIATOR = ${PROJECT_DIR}/cmd/logMediator
SCRIBE = ${PROJECT_DIR}/scribe
API = ${PROJECT_DIR}/api
CERT = ${PROJECT_DIR}/certs

GRPC_JAVA_PLUGIN = /home/romanos/go/bin/protoc-gen-grpc-java

dirs:
	echo PROJECT_NAME = ${PROJECT}
	echo GITHUB = ${GITHUB}
	echo PROJECT_DIR = ${PROJECT_DIR}
	echo LOG_SCRIBE_CMD = ${LOG_SCRIBE_CMD}
	echo SCRIBE = ${SCRIBE}

test:
	cd ${SCRIBE} && go test -v && cd ${PROJECT_DIR}

clean:
	if test -f  ${LOG_SCRIBE_CMD}/logScribe; \
	then rm -rf ${LOG_SCRIBE_CMD}/logScribe; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -f  ${MEDIATOR}/logMediator; \
	then rm -rf ${MEDIATOR}/logMediator; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -f  ${LOG_SCRIBE_CMD}/logScribe.exe; \
	then rm -rf ${LOG_SCRIBE_CMD}/logScribe.exe; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -f  ${SCRIBE}/file.txt; \
	then rm -rf ${SCRIBE}/file.txt; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -d  ${SCRIBE}/testdata; \
	then rm -rf ${SCRIBE}/testdata; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -d  ${SCRIBE}/noPath; \
	then rm -rf ${SCRIBE}/noPath; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	#clear
	echo "everything is clean"

buildScribe:
	cd ${LOG_SCRIBE_CMD} && go build && cd ${PROJECT_DIR}

buildMediator:
	cd ${MEDIATOR} && go build && cd ${PROJECT_DIR}

build: clean test buildScribe buildMediator

runScribe:
	${LOG_SCRIBE_CMD}/logScribe -path logs -pprof

runMediator:
	${MEDIATOR}/logMediator -pprof -pport 1122

secRun:
	${LOG_SCRIBE_CMD}/logScribe -path logs -pprof -crt ${CERT}/server.crt \
		-pk ${CERT}/server.key -ca ${CERT}/CertAuth.crt

all: build runScribe

clearLogs:
	rm -rf ${PROJECT_DIR}/logs
	mkdir ${PROJECT_DIR}/logs

# compiling .proto files
protoGo:
	cd ${API} && protoc --go_out=plugins=grpc:. *.proto && cd ${PROJECT_DIR}

protoJava:
	cd ${API} && protoc --grpc-java_out=. \
		--plugin=protoc-gen-grpc=${GRPC_JAVA_PLUGIN} \
		--java_out=. *.proto && cd ${PROJECT_DIR}


# dummy certificates for client, server, ans certificate authority
cleanCert:
	if test -d  ${CERT}; \
	then rm -rf ${CERT}; \
	fi
	mkdir ${CERT}

certAuth: cleanCert
	certstrap --depot-path ${CERT} init --cn "CertAuth"

cstrap:
	certstrap --depot-path ${CERT} request-cert -ip 127.0.0.1
	certstrap --depot-path ${CERT} sign 127.0.0.1 --CA CertAuth
	mv ${CERT}/127.0.0.1.crt ${CERT}/server.crt
	mv ${CERT}/127.0.0.1.key ${CERT}/server.key
	mv ${CERT}/127.0.0.1.csr ${CERT}/server.csr
	certstrap --depot-path ${CERT} request-cert -ip 127.0.0.1
	certstrap --depot-path ${CERT} sign 127.0.0.1 --CA CertAuth
	mv ${CERT}/127.0.0.1.crt ${CERT}/client.crt
	mv ${CERT}/127.0.0.1.key ${CERT}/client.key
	mv ${CERT}/127.0.0.1.csr ${CERT}/client.csr

# install dependencies
deps:
	go get google.golang.org/grpc
	go get google.golang.org/grpc/credentials
	go get google.golang.org/grpc/reflection
	go get github.com/golang/protobuf/proto
	go get github.com/rs/xid


# docker, is not ready yet
dockerBuildScribe:
	docker build -f cmd/logScribe/Dockerfile -t romanos/scribe cmd/logScribe/

dockerRunScribe:
	docker run -it --rm -v ${PWD}/logs:/logs --name scribe-service romanos/scribe -p 8080:8080 -p 1000:1111
