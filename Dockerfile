# Builds a docker image according to this Dockerfile by building the sysl binary
# on Go 1.13 (by default) and then copying it
# to a debian:buster-slim image and setting up the entrypoint.
#
# The produced image is published to https://hub.docker.com/r/anzbank/sysl

ARG go_ver=1.13
ARG deb_ver=10.3-slim

FROM golang:${go_ver} as builder

ARG sysl_tag

RUN test -n "${sysl_tag}" || ( echo Plese set sysl_tag build-arg && exit 1 )

WORKDIR /sysl

# COPY . .

RUN git clone --branch ${sysl_tag} --single-branch --depth 1 https://github.com/anz-bank/sysl.git . && \
    go install ./cmd/sysl

FROM debian:${deb_ver} as runner

ARG DEBIAN_FRONTEND=noninteractive

# Hack to fix ln problem when installing java headless
RUN mkdir -p /usr/share/man/man1 ; \
    apt-get update && \
    apt-get install -y plantuml

COPY --from=builder  /go/bin/sysl /

WORKDIR /work

ENTRYPOINT ["/sysl"]

CMD ["help"]
