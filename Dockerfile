FROM scratch

COPY sysl /

ENTRYPOINT ["/sysl"]
