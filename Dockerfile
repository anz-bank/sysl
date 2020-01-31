# Goreleaser builds a docker image according to this Dockerfile by copying the sysl binary 
# to a golang:1.13-alpine image and setting up the entrypoint.
#
# The produced image is published to https://hub.docker.com/r/anzbank/sysl

FROM golang:1.13-alpine

COPY sysl /

ENTRYPOINT ["/sysl"]
