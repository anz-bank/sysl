# Builds a docker image according to this Dockerfile by building the sysl binary
# on Go 1.13 (by default) using the current workspace and then copying the bianry
# to a image and setting up the entrypoint.
#
# The produced image is published to https://hub.docker.com/r/anzbank/sysl

ARG go_ver=1.13
ARG alpine_ver=3.11

FROM golang:${go_ver}-alpine${alpine_ver} as builder

RUN apk --no-cache add git

WORKDIR /sysl

COPY . .

RUN go install ./cmd/sysl

FROM golang:${go_ver}-alpine${alpine_ver} as runner

COPY --from=builder /go/bin/sysl /

WORKDIR /work

ENTRYPOINT ["/sysl"]

CMD ["help"]
