FROM ubuntu:18.04

RUN apt-get -qq update -y --fix-missing && \
  apt-get -qq upgrade -y --fix-missing && \
  apt-get -qq install -y python2.7 python-pip
  
WORKDIR /sysl

ADD dist/sysl-*.whl /sysl
ADD gosysl/gosysl-linux /usr/bin/sysl2

RUN pip install sysl-*.whl
RUN ln -s /usr/bin/sysl2 /usr/bin/syslgen
