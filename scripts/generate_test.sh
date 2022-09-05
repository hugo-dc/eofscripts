#!/bin/bash

# Params:
# 1: returned bytecode data
# 2: deploying eof data

echo "==============================================================="
terminating_opcodes="stop return revert invalid selfdestruct"
for to in $terminating_opcodes ; do
  retbc=$(./build/eof_gen "$(cat $1.txt)" "$(cat $to.txt)" | tail -n1 )

  echo "ASM:               $(cat asm_$to.txt)"
  echo "Returned bytecode: $retbc"

  yul=$(./build/yulreturn $retbc)
  echo "Yul code:          $yul"
  cyul=$(yul_comp "$yul" | tail -n1)
  echo "Compiled yul:      $cyul"
  eofyul=$(./build/eof_gen "$(cat $2.txt)" "$cyul" | tail -n1)
  echo "EOF Yul:           $eofyul"

  echo "-----------------------------------------------------------------"
done
