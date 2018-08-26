FROM alpine:latest

LABEL Romanos Trechlis "r.trechlis@gmail.com"

COPY go-scribe /usr/local/bin/go-scribe

WORKDIR /logs

ENTRYPOINT ["go-scribe"]

ENV AGENT_PORT 8080
ENV PPROF_PORT 1111
ENV RPC_PORT 4242


EXPOSE $AGENT_PORT $PPROF_PORT $RPC_PORT
