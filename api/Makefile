go:
	protoc --go_out=plugins=grpc:. *.proto

java:
	mkdir ../java/src/api
	protoc --java_out ../java/src --plugin=grpc *.proto
