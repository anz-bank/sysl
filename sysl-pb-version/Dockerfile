ARG ARRAI_VERSION="v0.220.0"
ARG ARRAI_IMAGE=anzbank/arrai
ARG SYSL_IMAGE=anzbank/sysl

FROM ${ARRAI_IMAGE}:${ARRAI_VERSION} AS stage

WORKDIR /sysl
COPY . .
RUN go get github.com/anz-bank/go-bindata/...
RUN make build
RUN arrai bundle -o config.arraiz sysl-pb-version/config_cli.arrai

FROM ${ARRAI_IMAGE}:${ARRAI_VERSION}

COPY --from=stage /sysl/dist/sysl /bin/

WORKDIR /scripts
COPY --from=stage /sysl/config.arraiz ./
COPY --from=stage /sysl/sysl-pb-version/version.sh ./

WORKDIR /work
ENTRYPOINT [ "/scripts/version.sh" ]
