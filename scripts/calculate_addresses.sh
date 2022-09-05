#!/bin/bash

cd ~/workspace/ethscripts/

i=0
ADDRESS=$1

# If nonce is provided we only calculate one address
if [ ! -z "$2" ]; then
  i=$2
  result=$(./build/create_address $ADDRESS $i)
  echo "${result:2}:  # $i"
  exit
fi


for i in {0..256}
do
  result=$(./build/create_address $ADDRESS $i)
  echo "${result:2}:  # PUSH$(( $i + 1 ))"
  echo "  nonce: 1"
  echo "  storage: {}"
  echo "  code: ''"
done
