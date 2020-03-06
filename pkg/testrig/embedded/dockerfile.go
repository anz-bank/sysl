package embedded

func GetDockerfileStub() string {
	return `
FROM golang:latest as builder

ARG GOPROXY
ENV GOPROXY=${GOPROXY}

WORKDIR /app

# allow caching on source rebuilds
COPY vendor ./vendor
COPY go.mod go.sum ./
RUN go mod download

# main build
COPY . ./
RUN go build -o main {{.OutputDir}}/{{.Service.Name}}/main.go



FROM scratch
COPY --from=builder /app/main .
CMD ["/main"]	
`
}
