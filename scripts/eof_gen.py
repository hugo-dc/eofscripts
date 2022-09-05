#!/bin/env python
import sys

data = sys.argv[1]
code = sys.argv[2]

print('code: ', code)
print('data: ', data)

code_len = len(code)
data_len = len(data)

if code_len % 2 != 0:
    print("Error: odd code size")
    exit(1)

if data_len % 2 != 0:
    print("Error: odd data size")
    exit(2)

code_len = int(code_len / 2)
data_len = int(data_len / 2)


code_len_hex = hex(code_len)
code_len_hex = code_len_hex[2:]

if len(code_len_hex) % 2 != 0:
    code_len_hex = '0' + code_len_hex

if len(code_len_hex) < 4:
    code_len_hex = '00' + code_len_hex

print('code length', code_len, code_len_hex)
code_section = '01' + code_len_hex
data_section = ''

data_len_hex = ''
if data_len > 0:
    data_len_hex = hex(data_len)
    data_len_hex = data_len_hex[2:]

    if len(data_len_hex) % 2 != 0:
        data_len_hex = '0' + data_len_hex


    if len(data_len_hex) < 4:
        data_len_hex = '00' + data_len_hex

    data_section = '02' + data_len_hex

print('data length', data_len, data_len_hex)

terminator = '00'

eof_code = 'ef0001' + code_section + data_section + terminator + code + data

print(eof_code)

