FROM golang:latest
#!/bin/bash
ENV WORKDIR $GOPATH/src/mongodb-on-demand
RUN mkdir -p $WORKDIR
RUN pwd
RUN go env
RUN ls
RUN uname -s
RUN uname -p
WORKDIR $WORKDIR

ADD . $WORKDIR/
RUN make test