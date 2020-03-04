# Builds a docker image according to this Dockerfile by building the sysl binary
# on Go 1.13 (by default) and then copying it
# to a debian:buster-slim image and setting up the entrypoint.
#
# The produced image is published to https://hub.docker.com/r/anzbank/sysl
#
# To use this Docker file run the following command
#    docker image build --build-arg sysl_tag=v0.7.1  -t my-sysl:0.7.1 .

ARG go_ver=1.13
ARG alpine_ver=3.11

FROM golang:${go_ver}-alpine${alpine_ver} as builder

ARG sysl_tag

RUN test -n "${sysl_tag}" || ( echo Plese set sysl_tag build-arg && exit 1 )

RUN apk --no-cache add git

WORKDIR /sysl

RUN git clone --branch ${sysl_tag} --single-branch --depth 1 https://github.com/anz-bank/sysl.git . && \
    go install ./cmd/sysl

FROM alpine:${alpine_ver} as runner

COPY --from=builder /go/bin/sysl /

WORKDIR /work

ENTRYPOINT ["/sysl"]

CMD ["help"]
