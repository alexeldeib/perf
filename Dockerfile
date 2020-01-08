# FROM ubuntu:disco-20191030
FROM ubuntu@sha256:550b0f99b69b0a1dec737f606ee294546466ff281036d4be6aff09bc1ac006fe as base
RUN apt update -y
RUN apt-get install -y git
RUN apt-get install -y nano 
RUN apt-get install -y jq 
RUN apt-get install -y fio 
RUN apt-get install -y iotop 
RUN apt-get install -y sysstat 
RUN apt-get install -y bpfcc-tools 
RUN apt-get install -y bpftrace 
RUN apt-get install -y blktrace 
RUN apt-get install -y linux-tools-common
RUN apt-get install -y magic-wormhole 
RUN apt-get install -y build-essential
RUN git clone https://github.com/brendangregg/FlameGraph.git /flamegraph

# prewarm module cache
FROM golang:1.13.4 as prewarmer

RUN mkdir /work
WORKDIR /work

ADD go.mod .
ADD go.sum .

RUN go mod download

# build extra tools
FROM prewarmer as builder

ADD . .

RUN go build -o aced .

FROM base as final
ENV PATH=$PATH:/work/az:/flamegraph
RUN mkdir /az

COPY --from=builder /work/aced ./az/aced