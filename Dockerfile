FROM golang:1.13-alpine

COPY sysl /

ENTRYPOINT ["/sysl"]
