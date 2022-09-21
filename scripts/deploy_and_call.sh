#!/bin/bash

SERVER="http://192.168.10.243"
RPC_SERVER="$SERVER:8545"
EXP_SERVER="$SERVER:4000"

bcode=$1

initcode=$(./build/returnevm $1)

./build/deploy $initcode

tx=""
while true; do
  echo "Fetching transactions..."

  lastBlock=$(curl --silent --location --request POST $RPC_SERVER --header 'Content-Type: application/json' --data-raw ' {
    "jsonrpc": "2.0",
    "method": "eth_getBlockByNumber",
    "params": ["latest", false],
    "id": 1
  }')

  echo "Block: $(echo $lastBlock | jq '.number') " 
  txs=$(echo $lastBlock | jq '.result.transactions')

  if [ "$txs" = "[]" ]; then
    echo "not found..."
  else
    tx=$(echo $txs | jq --raw-output '.[0]')
    break
  fi
  sleep 1
done

echo "Transaction found: $tx"

xdg-open $EXP_SERVER/tx/$tx

txReceipt=$(curl --silent --location --request POST $RPC_SERVER --header 'Content-Type: application/json' --data-raw " {
    \"jsonrpc\": \"2.0\",
    \"method\": \"eth_getTransactionReceipt\",
    \"params\": [\"$tx\"],
    \"id\": 1
  }")

contractAddress=$(echo $txReceipt | jq --raw-output '.result.contractAddress')
echo "Created contract: $contractAddress"

./build/call $contractAddress

xdg-open $EXP_SERVER/address/$contractAddress
