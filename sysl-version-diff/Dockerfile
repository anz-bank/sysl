ARG SYSL_PB_VERSION_IMAGE=anzbank/sysl-pb-version
ARG SYSL_PB_VERSION_VERSION="latest"

FROM ${SYSL_PB_VERSION_IMAGE}:${SYSL_PB_VERSION_VERSION} AS stage

WORKDIR /sysl
COPY . .
RUN arrai bundle -o version_diff.arraiz ./sysl-version-diff/version_diff_cli.arrai

FROM ${SYSL_PB_VERSION_IMAGE}:${SYSL_PB_VERSION_VERSION}
WORKDIR /scripts
COPY --from=stage /sysl/sysl-version-diff/*.sh ./
COPY --from=stage /sysl/version_diff.arraiz .

WORKDIR /work
ENTRYPOINT [ "/scripts/version-diff.sh" ]
