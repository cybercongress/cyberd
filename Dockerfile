###############################################################################
# Build cyber
###############################################################################
FROM ubuntu:18.04 as build_stage_cuda

ENV GO_VERSION 1.13.14
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '32617db984b18308f2b00279c763bff060d2739229cb8037217a49c9e691b46a'
ENV DAEMON_HOME /root/.cyberd
ENV CLIENT_HOME /root/.cyberdcli
ENV DAEMON_RESTART_AFTER_UPGRADE=on
ENV GAIA_HOME ${DAEMON_HOME}
ENV DAEMON_NAME cyberd
ENV BUILD_DIR /build
ENV COSMWASM_VER "0.7.2"
ENV PATH /usr/local/go/bin:/root/.cargo/bin:/root/.cyberd/scripts:$PATH


# Install required dev tools to compile cyberd
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget git ca-certificates

# Install golang
###############################################################################
RUN url="https://golang.org/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz" && \
    wget -O go.tgz "$url" && \
    echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tgz &&\
    rm go.tgz

ENV PATH="/usr/local/go/bin:$PATH"
RUN apt-get -y install --no-install-recommends \
    make gcc g++ \
    wget \
    curl \
    git \
    nvidia-cuda-toolkit \
&& go version

# Create appropriate folders layout
###############################################################################
 RUN mkdir -p /cyberd/upgrade_manager/genesis/bin \
 && mkdir -p /cyberd/upgrade_manager/upgrades/darwin/bin

# Compile cosmosd
###############################################################################
 RUN git clone https://github.com/regen-network/cosmosd.git $BUILD_DIR \
 && cd $BUILD_DIR \
 && go build \
 && cp cosmosd /usr/bin/cosmosd \
 && mv cosmosd /cyberd \
 && chmod +x /cyberd/cosmosd \
 && rm -fR $BUILD_DIR/* && rm -fR $BUILD_DIR/.*[a-z]

COPY . /sources
WORKDIR /sources/x/rank/cuda
RUN make build
RUN cp ./build/libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/

# Compile cyberd for genesis version
###############################################################################
WORKDIR /sources
RUN git checkout v0.1.6.2
RUN make build \
&& ./build/cyberd version \
&& cp ./build/cyberd /cyberd/upgrade_manager/genesis/bin/cyberd

# Compile cyberd for darwin upgrade version
###############################################################################
WORKDIR /sources
RUN git checkout v0.1.6.5
WORKDIR /sources/x/rank/cuda
RUN make 
RUN cp ./build/libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/
WORKDIR /sources
RUN make build \
&& ./build/cyberd version \
&& cp  ./build/cyberd /cyberd/upgrade_manager/upgrades/darwin/bin/cyberd

###############################################################################
# Build wasmvm
###############################################################################

RUN curl https://sh.rustup.rs -sSf | sh -s -- -y \
 && wget --quiet https://github.com/CosmWasm/wasmvm/archive/refs/tags/v${COSMWASM_VER}.tar.gz -P /tmp \
 && tar xzf /tmp/v${COSMWASM_VER}.tar.gz -C $BUILD_DIR \
 && cd $BUILD_DIR/wasmvm-${COSMWASM_VER}/ && make build \
 && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libgo_cosmwasm.so /usr/lib/ \
 && cp $BUILD_DIR/wasmvm-${COSMWASM_VER}/api/libgo_cosmwasm.dylib /usr/lib/

###############################################################################
# Create runtime cyber image
###############################################################################
FROM nvidia/cuda:10.2-base

ENV DAEMON_HOME /root/.cyberd
ENV CLIENT_HOME /root/.cyberdcli
ENV DAEMON_RESTART_AFTER_UPGRADE=on
ENV GAIA_HOME ${DAEMON_HOME}
ENV DAEMON_NAME cyberd

# Install useful dev tools
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget curl

# Download genesis file and links file from IPFS
###############################################################################
# To slow using ipget, currently we use gateway
# PUT needed CID_OF_GENESIS and CID_OF_CONFIG here
RUN wget -O /genesis.json https://ipfs.io/ipfs/QmZHpLc3H5RMXp3Z4LURNpKgNfXd3NZ8pZLYbjNFPL6T5n
RUN wget -O /config.toml https://ipfs.io/ipfs/QmSEPs57PyaK5envPyJ16jQZ9dwfhDePz6fmG4WaLWFVts

WORKDIR /

# Copy compiled kernel and binaries for current bin version
###############################################################################
COPY --from=build_stage_cuda /cyberd /cyberd

COPY --from=build_stage_cuda /sources/build/cyberd /usr/bin/cyberd
COPY --from=build_stage_cuda /sources/build/cyberdcli /usr/bin/cyberdcli

COPY --from=build_stage_cuda /usr/lib/cbdrank.h /usr/lib/cbdrank.h
COPY --from=build_stage_cuda /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so

COPY --from=build_stage_cuda /usr/lib/libgo_cosmwasm.so /usr/lib/libgo_cosmwasm.so
COPY --from=build_stage_cuda /usr/lib/libgo_cosmwasm.dylib /usr/lib/libgo_cosmwasm.dylib

# Copy startup scripts
###############################################################################

COPY start_script.sh /start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x start_script.sh
RUN chmod +x /entrypoint.sh

# Cleanup for runtime container
###############################################################################
#RUN cd / && rm -fR $BUILD_DIR $HOME/.rustup $HOME/.cargo $HOME/.cache $HOME/go 
#&& apt-get remove -y --auto-remove \
#ca-certificates \
#make gcc g++ \
#curl \
#git \
#nvidia-cuda-toolkit \
#; \
#rm -rf /var/lib/apt/lists/*; \
#rm -fR /tmp/*;

# Start
###############################################################################
EXPOSE 26656 26657 1317
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]