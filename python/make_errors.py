# Script to convert ghidra_error_messages.txt into a string array embedded in a .go
# file

import re

if __name__ == '__main__':
    # open file and put each line as an item in a list
    with open('ghidra_error_messages.txt', 'r') as f:
        lines = f.readlines()
    # iterate the list and extract strings
    r = re.compile("\"(.*?)\"")
    strings = []
    for line in lines:
        results = r.findall(line)
        if len(results) > 0:
            strings.append(results[0])
    # write errors.go
    with open('errors.go', 'w') as f:
        f.write('package main\n\n')
        f.write('// Error codes defined here\n')
        f.write('const (\n')
        counter = 2000
        for string in strings:
            variable = string.title().replace('/', '').replace('\'', '').replace('(', '').replace(')', '').replace(':', '').replace(' ', '')
            f.write(f'    Er{variable} = {counter}\n')
            counter += 1
        f.write(')\n\n')
        f.write('// ErrorMessages returns a map of error codes and their messages\n')
        f.write('func ErrorMessages() map[int]string {\n')
        f.write('    return map[string]int{\n')
        for string in strings:
            variable = string.title().replace('/', '').replace('\'', '').replace('(', '').replace(')', '').replace(':', '').replace(' ', '')
            f.write(f'        Er{variable}: \"{string}\"\n')
        f.write('    }\n')
        f.write('}\n')