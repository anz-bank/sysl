# Builds a docker image by building the sysl binary
# on Go 1.22 (by default) using the current workspace and then copying the binary
# to an image and setting up the entrypoint.
#
# The produced image is published to https://hub.docker.com/r/anzbank/sysl

ARG go_ver=1.22
ARG alpine_ver=3.18

ARG DOCKER_BASE=golang:${go_ver}-alpine${alpine_ver}
FROM ${DOCKER_BASE} AS builder

RUN apk --no-cache add git make

WORKDIR /sysl

COPY . .

RUN make build

FROM ${DOCKER_BASE} AS runner

COPY --from=builder /sysl/dist/sysl /

ENTRYPOINT ["/sysl"]

CMD ["help"]
