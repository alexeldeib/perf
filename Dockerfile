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
RUN apt-get install -y linux-headers-5.0.0-1020-azure linux-headers-5.0.0-1018-azure
RUN apt-get install -y linux-tools-5.0.0-1020-azure linux-tools-5.0.0-1018-azure
RUN apt-get install -y magic-wormhole 
RUN apt-get install -y build-essential
RUN git clone https://github.com/brendangregg/FlameGraph.git /flamegraph
RUN wget --quiet -O - https://deb.nodesource.com/gpgkey/nodesource.gpg.key | apt-key add - \
 && VERSION=node_12.x \
 && DISTRO="$(lsb_release -s -c)" \
 && echo "deb https://deb.nodesource.com/$VERSION $DISTRO main" | tee /etc/apt/sources.list.d/nodesource.list \
 && echo "deb-src https://deb.nodesource.com/$VERSION $DISTRO main" | tee -a /etc/apt/sources.list.d/nodesource.list \
 && apt update -y \
 && apt install -y nodejs
RUN npm i -g 0x
RUN npm i -g stackvis

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

# bpfbase v0.0.1-1018 sha256:cd42b78e598a08f5203db0139025676b690077a424c7863fae1ed5680f435b87
# bpfbase v0.0.1-1020 sha256:4ae9eb789d91631480440c8e35185f5768cc15c9d3851a041a02d2b5a98650c9
FROM base as final
ENV PATH=$PATH:/work/az:/flamegraph
RUN mkdir /az

COPY --from=builder /work/aced ./az/aced