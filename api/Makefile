go:
	protoc --go_out=plugins=grpc:. *.proto

java:
	protoc --grpc-java_out=. --plugin=protoc-gen-grpc=/home/romanos/GO/bin/protoc-gen-grpc-java --java_out=. *.proto
