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
    llvm-7-dev \
    llvm-7-runtime \
    libclang-7-dev \
    clang-7 \
    python \
    zlib1g-dev \
    libelf-dev

RUN git clone https://github.com/iovisor/bcc.git \
    && mkdir bcc/build; cd bcc/build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local .. \
    && make \
    && make install

RUN git clone https://github.com/iovisor/bpftrace.git \
    && mkdir bpftrace/build; cd bpftrace/build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local ..  \
    && make \
    && make install

FROM ubuntu:${UBUNTU_TAG}
RUN apt update -y && apt install -y python3
COPY --from=base /usr/local/bcc /usr/
COPY --from=base /usr/local/bpftrace /usr/
