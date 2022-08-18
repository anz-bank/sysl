# Builds a docker image by building the sysl binary
# on Go 1.19 (by default) using the current workspace and then copying the binary
# to an image and setting up the entrypoint.
#
# The produced image is published to https://hub.docker.com/r/anzbank/sysl

ARG go_ver=1.19
ARG alpine_ver=3.16

FROM golang:${go_ver}-alpine${alpine_ver} as builder

RUN apk --no-cache add git make

WORKDIR /sysl

COPY . .

RUN make build

FROM golang:${go_ver}-alpine${alpine_ver} as runner

COPY --from=builder /sysl/dist/sysl /

ENTRYPOINT ["/sysl"]

CMD ["help"]
