import sys

for line in sys.stdin:
    line = line.strip()
    temperature = int(line[87:92])
    if ((temperature != 9999)):
        print('%s\t%d' % (line[15:19], int(line[87:92])))