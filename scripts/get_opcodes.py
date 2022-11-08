import sys

def get_valid_opcodes():
    vo_ranges = [
        (0x00, 0x0b), 
        (0x10, 0x1d), 
        (0x20, 0x20), 
        (0x30, 0x3f), 
        (0x40, 0x48), 
        (0x50, 0x5b), 
        (0x60, 0x6f), 
        (0x70, 0x7f), 
        (0x80, 0x8f), 
        (0x90, 0x9f), 
        (0xa0, 0xa4), 
        (0xf0, 0xf5), 
        (0xfa, 0xfa), 
        (0xfd, 0xfd), 
        (0xfe, 0xfe), 
        (0xff, 0xff)]

    valid_opcodes = []
    for r in vo_ranges:
        for i in range(r[0], r[1] + 1):
            valid_opcodes.append(i)
    return valid_opcodes

def get_terminating_opcodes():
    terminating_opcodes = [0x00, 0xf3, 0xfd, 0xfe, 0xff]
    return terminating_opcodes

def get_non_terminating_opcodes():
    valid_opcodes = get_valid_opcodes()
    terminating_opcodes = get_terminating_opcodes()
    non_term_opcodes = [x for x in valid_opcodes if x not in terminating_opcodes]
    return non_term_opcodes

def print_opcodes(op_list):
    for op in op_list:
        op_hex = hex(op)
        if len(op_hex) % 2 != 0:
            op_hex = '0x0' + op_hex[2:] 
        print(op_hex)

option="all"
if len(sys.argv) == 2:
    option=sys.argv[1]

if option == '--all' or option == '--valid':
    print_opcodes( get_valid_opcodes() )

if option == '--terminating-opcodes' or option == '-to':
    print_opcodes( get_terminating_opcodes () )

if option == '--non-terminating-opcodes' or option == '-nto':
    print_opcodes( get_non_terminating_opcodes() )

print
