ARG ARRAI_VERSION="v0.252.0"
ARG ARRAI_IMAGE=anzbank/arrai

FROM ${ARRAI_IMAGE}:${ARRAI_VERSION} AS stage

RUN apk add --no-cache make

WORKDIR /sysl
COPY . .
RUN make build
RUN arrai bundle -o version_diff.arraiz ./sysl-version-diff/version_diff_cli.arrai

FROM ${ARRAI_IMAGE}:${ARRAI_VERSION}
COPY --from=stage /sysl/dist/sysl /bin/

# copy sysl.pb for newNormalize
WORKDIR /sysl
COPY --from=stage /sysl/pkg/sysl/sysl.pb .

WORKDIR /scripts
COPY --from=stage /sysl/sysl-version-diff/*.sh ./
COPY --from=stage /sysl/version_diff.arraiz .

WORKDIR /workdir
ENTRYPOINT [ "/scripts/version-diff.sh" ]
