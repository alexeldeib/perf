ARG UBUNTU_TAG=16.04

FROM ubuntu:${UBUNTU_TAG} as base

ENV DEBIAN_FRONTEND noninteractive

WORKDIR /work

RUN apt update -y && apt -y install bison \
    build-essential \
    cmake \
    flex \
    git \
    libedit-dev \
    libllvm6.0 \
    llvm-6.0-dev \
    libclang-6.0-dev \
    python \
    zlib1g-dev \
    libelf-dev

RUN git clone https://github.com/iovisor/bcc.git \
    && mkdir bcc/build; cd bcc/build \
    && cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local/bcc \
    && make \
    && make install

# RUN git clone https://github.com/iovisor/bpftrace.git \
#     && mkdir bpftrace/build; cd bpftrace/build \
#     && cmake .. -DCMAKE_INSTALL_PREFIX=/usr/local/bpftrace \
#     && make \
#     && make install

FROM ubuntu:${UBUNTU_TAG}
COPY --from=base /usr/bin/python /usr/bin/python
COPY --from=base /usr/local/bcc /usr/
RUN apt update -y && apt install -y python
