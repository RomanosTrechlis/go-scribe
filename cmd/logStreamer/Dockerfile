FROM golang:alpine

# add git support to golang alpine
RUN apk add --update --no-cache git

COPY . /go/src/github.com/RomanosTrechlis/logStreamer/cmd/logStreamer

WORKDIR /romanos

RUN go install github.com/RomanosTrechlis/logStreamer/cmd/logStreamer

ENTRYPOINT /go/bin/logStreamer -path ./logs -pprof

EXPOSE 8080 1111
