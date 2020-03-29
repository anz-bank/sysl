package testrig

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
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main {{.OutputDir}}/{{.Service.Name}}/main.go



FROM scratch
COPY --from=builder /app/main /main
CMD ["/main"]	
`
}
