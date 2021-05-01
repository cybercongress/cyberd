#!/bin/bash
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
DELEGATOR='Delegator Address'
VALIDATOR='Operator Address'
PASWD='Keyring Password'
DELAY=7200 #in secs
ACC_NAME=Account Name

for (( ;; )); do
	BAL=$(cyberdcli query account ${DELEGATOR} --chain-id euler-6);
        echo -e "BALANCE: ${GREEN}${BAL}${NC} eul\n"
	echo -e "Claim rewards\n"
	echo -e "${PASWD}\n${PASWD}\n" | cyberdcli tx distribution withdraw-rewards ${VALIDATOR} --chain-id euler-6 --from ${ACC_NAME} --commission -y 
	for (( timer=10; timer>0; timer-- ))
	do
		printf "* sleep for ${RED}%02d${NC} sec\r" $timer
		sleep 1
	done
	BAL=$(cyberdcli -o json query account ${DELEGATOR} --chain-id euler-6 | jq -r '.value.coins[].amount');
	echo -e "BALANCE: ${GREEN}${BAL}${NC} eul\n"
	echo -e "Stake ALL\n"
	echo -e "${PASWD}\n${PASWD}\n" | cyberdcli tx staking delegate ${VALIDATOR} ${BAL}eul --chain-id euler-6 --from ${ACC_NAME} -y -o json | jq .  
	for (( timer=${DELAY}; timer>0; timer-- ))
	do
		printf "* sleep for ${RED}%02d${NC} sec\r" $timer
		sleep 1
	done
done
