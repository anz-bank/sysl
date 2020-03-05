# The produced image is published to https://hub.docker.com/r/anzbank/sysl

FROM golang:1.13.8-alpine3.11 as builder

# Allow custom GOPROXY if in a network-constrained environment (corporate)
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

# Allow custom Alpine archive if in a network-constrained evironment (corporate)
ARG ALPINE_ARCHIVE
RUN [[ -n "${ALPINE_ARCHIVE}" ]] && \
    sed -i -e "s|http://dl-cdn.alpinelinux.org|${ALPINE_ARCHIVE}|" /etc/apk/repositories
RUN apk add --no-cache git make

# /src is outside GOPATH so we don't need to set GO111MODULE=on
WORKDIR /src
COPY go.mod go.sum Makefile ./
COPY .git .git
COPY cmd cmd
COPY pkg pkg
COPY vendor vendor
RUN make build

# We use golang as the base image because sysl uses the `go` command to manage
# sysl modules.
FROM golang:1.13.8-alpine3.11
COPY --from=builder /src/dist/sysl /sysl
ENTRYPOINT ["/sysl"]
CMD ["help"]
