FROM alpine

RUN apk add --no-cache python py-pip
WORKDIR /sysl
ADD dist/sysl-*.whl /sysl
RUN pip install sysl-*.whl

CMD ["sysl", "-h"]
