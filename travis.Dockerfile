# From https://hub.docker.com/r/travisci/ci-opal/tags
FROM travisci/ci-opal:packer-1545391454-cb6d370

RUN apt-get -y update && \
    curl -sL https://deb.nodesource.com/setup_11.x | bash - && \
    apt-get -y install \
        build-essential \
        golang \
        gradle \
        nodejs \
        && \
    sudo apt-get clean && \
    pip install flake8==3.5.0

ENV GOPATH=/go \
    PATH=/go/bin:$PATH \
    TRAVIS_BUILD_DIR=/go/src/github.com/anz-bank/sysl \
    SYSL_PLANTUML=http://www.plantuml.com/plantuml

RUN mkdir -p $TRAVIS_BUILD_DIR
WORKDIR $TRAVIS_BUILD_DIR

RUN go get -v \
    github.com/antlr/antlr4/runtime/Go/antlr \
    github.com/golang/protobuf/proto         \
    github.com/golang/protobuf/ptypes/struct \
    github.com/golang/protobuf/jsonpb

ENV NPM_AUTH_TOKEN=SOME-RANDOM-KEY
RUN mkdir -p sysl2/sysl/sysl_js
COPY ./sysl2/sysl/sysl_js $TRAVIS_BUILD_DIR/sysl2/sysl/sysl_js
RUN npm install --prefix sysl2/sysl/sysl_js

COPY . .
RUN pip install . pytest
