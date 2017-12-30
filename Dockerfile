FROM alpine:latest

LABEL Romanos Trechlis  "r.trechlis@gmail.com"

COPY go-scribe /usr/local/bin/go-scribe

WORKDIR /logs

ENTRYPOINT ["go-scribe", "docker-agent"]

ENV AGENT_PORT 8080
ENV PPROF_PORT 1111
ENV PATH /usr/local/bin/:$PATH


EXPOSE $AGENT_PORT $PPROF_PORT
