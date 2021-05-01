#!/bin/sh

if [ ! -f "/root/.cyberd/config/node_key.json" ]
then
  mkdir /root/.cyberd/config/
  cyberd init ${NODE_MONIKER}
  cp /genesis.json /root/.cyberd/config/
  cp /config.toml  /root/.cyberd/config/
  cp /cyberd/cosmosd  /root/.cyberd
  cp -r /cyberd/upgrade_manager/  /root/.cyberd
#  cp /links        /root/.cyberd/config/
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