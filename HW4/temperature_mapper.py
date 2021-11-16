import sys

for line in sys.stdin:
    line = line.strip()
    temperature = int(line[87:92])
    q = int(line[92])
    if ((temperature != 9999)) and (q == 0 or q == 1 or q == 4 or q == 5 or q == 9) == True:
        print('%s\t%d' % (line[15:23], int(line[87:92])))
