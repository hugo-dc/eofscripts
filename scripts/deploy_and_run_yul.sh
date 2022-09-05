#!/bin/bash

nonce=$(cat ./.nonce)
./deploy_yul "$1"

