package test_rig

func GetDokerfileStub() string {
	return `FROM scratch
COPY ./main /main
CMD ["/main"]
`
}

func GetDockerfileStub2() string {
	return `
FROM golang:alpine as builder
ENV http_proxy=http://docker.for.mac.host.internal:3128
ENV https_proxy=http://docker.for.mac.host.internal:3128

# certificates needed for golang downloads
RUN apk --no-cache add ca-certificates
COPY *.cer /usr/local/share/ca-certificates/
RUN update-ca-certificates

WORKDIR /app

# dependencies as explicit step to allow caching
COPY go.mod go.sum ./
COPY vendor ./vendor
RUN go mod download

# main build
COPY . ./
RUN GOOS=linux go build -o main cmd/dbfront/main.go



FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main .
# ENTRYPOINT [ "/main" ]
CMD ["/main"]
`
}
