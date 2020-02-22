FROM nvidia/cuda:10.0-devel-ubuntu18.04 as build_stage

ENV GO_VERSION 1.13.1
ENV GO_ARCH 'linux-amd64'
ENV GO_BIN_SHA '94f874037b82ea5353f4061e543681a0e79657f787437974214629af8407d124'


#  Install required dev tools to compile cyberd
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget git


#  Install golang
###############################################################################
RUN url="https://golang.org/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz" && \
	wget -O go.tgz "$url" && \
	echo "${GO_BIN_SHA} *go.tgz" | sha256sum -c - && \
	tar -C /usr/local -xzf go.tgz &&\
	rm go.tgz

ENV PATH="/usr/local/go/bin:$PATH"
RUN go version && nvcc --version


# Compile cuda kernel
###############################################################################
COPY . /sources
WORKDIR /sources/x/rank/cuda
RUN make build
RUN cp ./build/libcbdrank.so /usr/lib/ && cp cbdrank.h /usr/lib/


# Compile cyberd
###############################################################################
WORKDIR /sources
RUN make build


###############################################################################
# Create base image
###############################################################################
FROM nvidia/cuda:10.0-runtime-ubuntu18.04

#  Install useful dev tools
###############################################################################
RUN apt-get update && apt-get install -y --no-install-recommends wget curl

#  Download genesis file and links file from IPFS
###############################################################################
# To slow using ipget, currently we use gateway
RUN wget -O /genesis.json https://ipfs.io/ipfs/QmTTBz7svGkrnigzBhUvnLKExY6LfYQrM3jXzWUFLeM3Xo
RUN wget -O /config.toml https://ipfs.io/ipfs/QmYDtVn8xJiR69fLDF2d384aSkSNzwKMrAwz5saa4HUuaj

WORKDIR /

#  Copy compiled kernel and binaries
###############################################################################
COPY --from=build_stage /sources/build/cyberd /usr/bin/cyberd
COPY --from=build_stage /sources/build/cyberdcli /usr/bin/cyberdcli

COPY --from=build_stage /usr/lib/cbdrank.h /usr/lib/cbdrank.h
COPY --from=build_stage /usr/lib/libcbdrank.so /usr/lib/libcbdrank.so

COPY --from=build_stage /sources/libgo_cosmwasm.so /usr/bin/libgo_cosmwasm.so
COPY --from=build_stage /sources/libgo_cosmwasm.so /usr/lib/libgo_cosmwasm.so

#  Copy startup scripts
###############################################################################

COPY start_script.sh start_script.sh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x start_script.sh
RUN chmod +x /entrypoint.sh


#  Start
###############################################################################
EXPOSE 26656 26657
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./start_script.sh"]
