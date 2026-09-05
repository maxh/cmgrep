#!/usr/bin/env python3

import random
import re
import socket
import sys

# With a node number, write just that machine's log. With no argument, use the
# hostname if we are on a VM, otherwise write all ten (for debugging purposes).
if len(sys.argv) > 1:
    nodes = [int(sys.argv[1])]
else:
    match = re.search(r"cs425-72(\d\d)", socket.gethostname())
    if match:
        nodes = [int(match.group(1))]
    else:
        nodes = range(1, 11)

for node in nodes:
    random.seed(425 + node)

    frequent = 0
    happensrarely = 0
    morerare = 0
    alsoquiterare = 0

    f = open("machine.%d.log" % node, "w")
    for i in range(20000):
        method = random.choice(["GET", "POST", "PUT", "DELETE"])
        path = random.choice(["/", "/average", "/basic", "/standard"])

        r = random.random()
        if r < 0.20:
            path = "/frequent"
            frequent += 1
        elif r < 0.21:
            path = "/happensrarely"
            happensrarely += 1
        elif node == 3 and r < 0.23:
            path = "/morerare"
            morerare += 1
        elif node in (1, 2, 7) and r < 0.24:
            path = "/alsoquiterare"
            alsoquiterare += 1

        ip = "10.0.%d.%d" % (node, random.randint(1, 254))
        status = random.choice([200, 200, 404, 500])
        f.write('%s - - "%s %s HTTP/1.1" %d\n' % (ip, method, path, status))
    f.close()

    print("machine.%d.log frequent=%d happensrarely=%d morerare=%d alsoquiterare=%d"
          % (node, frequent, happensrarely, morerare, alsoquiterare))
