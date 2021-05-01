#!/bin/sh

if [ -z "${COMPUTE_RANK_ON_GPU}" ]
then
  COMPUTE_GPU=true
else
  COMPUTE_GPU="${COMPUTE_RANK_ON_GPU}"
fi

if [ -z "${ALLOW_SEARCH}" ]
then
  ALLOW_SEARCH_FLAG=false
else
  ALLOW_SEARCH_FLAG="${ALLOW_SEARCH}"
fi

#export $(cat /.env) &
ulimit -n 4096 &
# Start the first process
/root/.cyberd/cosmosd start --compute-rank-on-gpu=${COMPUTE_GPU} --allow-search=${ALLOW_SEARCH_FLAG} --home /root/.cyberd  &
#status=$?
#if [ $status -ne 0 ]; then
#  echo "Failed to start cyberd: $status"
#  exit $status
#fi
#
## Start the second process
## PUT needed CHAIN_ID here
cyberdcli rest-server  --trust-node --chain-id=euler-6 --laddr=tcp://0.0.0.0:1317 --indent --home=/root/.cyberdcli
#status=$?
#if [ $status -ne 0 ]; then
#  echo "Failed to start cyberd light-client: $status"
#  exit $status
#fi

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container exits with an error
# if it detects that either of the processes has exited.
# Otherwise it loops forever, waking up every 60 seconds
#while sleep 60; do
#  ps aux |grep cyberd |grep -q -v grep
#  PROCESS_1_STATUS=$?
#  ps aux |grep cyberdcli |grep -q -v grep
#  PROCESS_2_STATUS=$?
#  # If the greps above find anything, they exit with 0 status
#  # If they are not both 0, then something is wrong
#  if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
#    echo "One of the processes has already exited."
#    exit 1
#  fi
#done