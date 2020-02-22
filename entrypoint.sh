#!/bin/sh

if [ ! -f "/root/.cyberd/config/node_key.json" ]
then
  mkdir /root/.cyberd/config/
  cyberd init
  cp /genesis.json /root/.cyberd/config/
  cp /config.toml  /root/.cyberd/config/
#  cp /links        /root/.cyberd/config/
  cp /libgo_cosmwasm.so /root/
fi

if [ ! -f "/root/.cyberd/config/genesis.json" ]
then
  cp /genesis.json /root/.cyberd/config/
fi

if [ ! -f "/root/.cyberd/config/config.toml" ]
then
  cp /config.toml /root/.cyberd/config/
fi

#if [ ! -f "/root/.cyberd/config/links" ]
#then
#  cp /links /root/.cyberd/config/
#fi

if [ "$2" = 'init' ]
then
  return 0
else
  exec "$@"
fi