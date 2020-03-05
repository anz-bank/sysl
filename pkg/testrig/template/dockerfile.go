package template

func GetDockerfileStub() string {
	return `
FROM golang:latest as builder
ENV http_proxy=http://docker.for.mac.host.internal:3128
ENV https_proxy=http://docker.for.mac.host.internal:3128

WORKDIR /app

# allow insecure download for ANZ root CA certificates, we need them for golang
RUN git config --global http.sslVerify false
RUN git clone --depth 1 https://github.service.anz/ocp/ocp-cacerts.git
RUN cp ocp-cacerts/global/*.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

# allow caching on source rebuilds
COPY go.mod go.sum ./
COPY vendor ./vendor
RUN go mod download

# main build
COPY . ./
RUN GOOS=linux go build -o main {{outputDir}}/{{name}}/main.go



FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main .
# ENTRYPOINT [ "/main" ]
CMD ["/main"]	
`
}
