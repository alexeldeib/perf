ARG distro_tag=18.04

FROM ubuntu:${distro_tag} as base

ARG distro=bionic
ARG bcc_ref="v0.12.0"
ARG LLVM_VERSION="8"
ENV DISTRO=$distro
ENV LLVM_VERSION=$LLVM_VERSION

RUN apt-get update && apt-get install -y curl gnupg &&\
    llvmRepository="\n\
deb http://apt.llvm.org/${DISTRO}/ llvm-toolchain-${DISTRO} main\n\
deb-src http://apt.llvm.org/${DISTRO}/ llvm-toolchain-${DISTRO} main\n\
deb http://apt.llvm.org/${DISTRO}/ llvm-toolchain-${DISTRO}-${LLVM_VERSION} main\n\
deb-src http://apt.llvm.org/${DISTRO}/ llvm-toolchain-${DISTRO}-${LLVM_VERSION} main\n" &&\
    echo $llvmRepository >> /etc/apt/sources.list && \
    curl -L https://apt.llvm.org/llvm-snapshot.gpg.key | apt-key add -

RUN apt-get update && apt-get install -y \
      bison \
      binutils-dev \
      cmake \
      flex \
      g++ \
      git \
      libelf-dev \
      zlib1g-dev \
      libiberty-dev \
      libbfd-dev \
      libedit-dev \
      clang-${LLVM_VERSION} \
      libclang-${LLVM_VERSION}-dev \
      libclang-common-${LLVM_VERSION}-dev \
      libclang1-${LLVM_VERSION} \
      llvm-${LLVM_VERSION} \
      llvm-${LLVM_VERSION}-dev \
      llvm-${LLVM_VERSION}-runtime \
      libllvm${LLVM_VERSION} \
      systemtap-sdt-dev \
      python3 \
      quilt

RUN git clone https://github.com/iovisor/bcc.git \
    && mkdir bcc/build; cd bcc/build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local .. \
    && make \
    && make install 

# RUN ls /usr/local/include/bcc -l

RUN git clone https://github.com/iovisor/bpftrace.git \
    && mkdir bpftrace/build; cd bpftrace/build \
    && cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local -DCMAKE_CXX_FLAGS="-include /usr/local/include/bcc/compat/linux/bpf.h -D__LINUX_BPF_H__" ..  \
    && make \
    && make install
