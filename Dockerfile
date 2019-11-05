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

# FROM ubuntu:disco-20191030
# FROM ubuntu@sha256:550b0f99b69b0a1dec737f606ee294546466ff281036d4be6aff09bc1ac006fe
# MOVED TO DOCKERFILE.base
# RUN apt update -y
# RUN apt install -y git
# RUN apt install -y nano 
# RUN apt install -y jq 
# RUN apt install -y fio 
# RUN apt install -y iotop 
# RUN apt install -y sysstat 
# RUN apt install -y bpfcc-tools 
# RUN apt install -y bpftrace 
# RUN apt install -y blktrace 
# RUN apt install -y linux-tools-common
# RUN apt install -y linux-headers-5.0.0-1018-azure linux-headers-5.0.0-1020-azure
# RUN apt install -y linux-tools-5.0.0-1018-azure linux-tools-5.0.0-1020-azure
# # for transmitting files outside
# RUN apt install -y magic-wormhole 
# RUN apt-install -y build-essential
# RUN git clone https://github.com/brendangregg/FlameGraph.git /flamegraph

# bpfbase v0.0.1-1018 sha256:08a065f485890c689732901c5db9be83802189c8875b7d26304f6d8c853a63d6
# bpfbase v0.0.1-1020 sha256:4ae9eb789d91631480440c8e35185f5768cc15c9d3851a041a02d2b5a98650c9
FROM alexeldeib/bpfbase@sha256:4ae9eb789d91631480440c8e35185f5768cc15c9d3851a041a02d2b5a98650c9
ENV PATH=$PATH:/work/az:/flamegraph
RUN mkdir /az

COPY --from=builder /work/aced ./az/aced