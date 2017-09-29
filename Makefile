PROJECT = logStreamer
GITHUB = /home/romanos/GO/src/github.com/RomanosTrechlis
PROJECT_DIR = ${GITHUB}/${PROJECT}
LOG_STREAMER_CMD = ${PROJECT_DIR}/cmd/${PROJECT}
STREAMER = ${PROJECT_DIR}/streamer
API = ${PROJECT_DIR}/api
CERT = ${PROJECT_DIR}/certs

GRPC_JAVA_PLUGIN = /home/romanos/GO/bin/protoc-gen-grpc-java

dirs:
	echo PROJECT_NAME = ${PROJECT}
	echo GITHUB = ${GITHUB}
	echo PROJECT_DIR = ${PROJECT_DIR}
	echo LOG_STREAMER_CMD = ${LOG_STREAMER_CMD}
	echo STREAMER = ${STREAMER}

test:
	cd ${STREAMER} && go test -v && cd ${PROJECT_DIR}

clean:
	if test -f  ${LOG_STREAMER_CMD}/logStreamer; \
	then rm -rf ${LOG_STREAMER_CMD}/logStreamer; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -f  ${LOG_STREAMER_CMD}/logStreamer.exe; \
	then rm -rf ${LOG_STREAMER_CMD}/logStreamer.exe; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -f  ${STREAMER}/file.txt; \
	then rm -rf ${STREAMER}/file.txt; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -d  ${STREAMER}/testdata; \
	then rm -rf ${STREAMER}/testdata; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	if test -d  ${STREAMER}/noPath; \
	then rm -rf ${STREAMER}/noPath; \
	else echo "file doesn't exist. nothing to do"; \
	fi
	#clear
	echo "everything is clean"

build: clean test
	cd ${LOG_STREAMER_CMD} && go build && cd ${PROJECT_DIR}

run:
	${LOG_STREAMER_CMD}/logStreamer -path logs -pprof

secRun:
	${LOG_STREAMER_CMD}/logStreamer -path logs -pprof -crt ${CERT}/server.crt \
		-pk ${CERT}/server.key -ca ${CERT}/CertAuth.crt

all: build run

clearLogs:
	rm -rf ${PROJECT_DIR}/logs
	mkdir ${PROJECT_DIR}/logs

protoGo:
	cd ${API} && protoc --go_out=plugins=grpc:. *.proto && cd ${PROJECT_DIR}

protoJava:
	cd ${API} && protoc --grpc-java_out=. \
		--plugin=protoc-gen-grpc=${GRPC_JAVA_PLUGIN} \
		--java_out=. *.proto && cd ${PROJECT_DIR}

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
