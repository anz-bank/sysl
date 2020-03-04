# Goreleaser builds a docker image according to this Dockerfile by copying the sysl binary 
# to a golang:1.13-alpine image and setting up the entrypoint.
#
# The produced image is published to https://hub.docker.com/r/anzbank/sysl
ARG go_ver=1.14

FROM golang:${go_ver}-buster as builder

ARG sysl_ver=v0.7.1

WORKDIR /work

RUN git clone --branch ${sysl_ver} --single-branch --depth 1 https://github.com/anz-bank/sysl.git && \
    cd sysl && \
    go get ./... && \
    go install ./cmd/sysl

# RUN find / type f -name sysl -ls

FROM debian:buster-slim as runner

COPY --from=builder  /go/bin/sysl /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/sysl"]

CMD ["help"]