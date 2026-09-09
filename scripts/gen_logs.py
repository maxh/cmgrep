#!/usr/bin/env python3

import random
import re
import socket
import argparse

# With a node number, write just that machine's log. With no argument, use the
# hostname if we are on a VM, otherwise write all ten (for debugging purposes).
parser = argparse.ArgumentParser()
parser.add_argument("node", nargs="?", type=int)
parser.add_argument("--size-mb", type=int)
options = parser.parse_args()

if options.node is not None and not 1 <= options.node <= 10:
    parser.error("node must be between 1 and 10")

if options.size_mb is not None and options.size_mb <= 0:
    parser.error("size-mb must be positive")

target_bytes = None
if options.size_mb is not None:
    target_bytes = options.size_mb * 1_000_000

if options.node is not None:
    nodes = [options.node]
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

    f = open("machine.%d.log" % node, "w", encoding="ascii", newline="\n")
    bytes_written = 0
    lines_written = 0

    if target_bytes is not None:
        marker = '10.0.%d.1 - - "GET /rare HTTP/1.1" 200\n' % node
        f.write(marker)
        bytes_written += len(marker.encode("ascii"))
        lines_written += 1

    while True:
        if target_bytes is None and lines_written >= 20000:
            break
        if target_bytes is not None and bytes_written >= target_bytes:
            break
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
        line = '%s - - "%s %s HTTP/1.1" %d\n' % (ip, method, path, status)
        f.write(line)
        bytes_written += len(line.encode("ascii"))
        lines_written += 1
    f.close()

    print("machine.%d.log frequent=%d happensrarely=%d morerare=%d alsoquiterare=%d"
          % (node, frequent, happensrarely, morerare, alsoquiterare))
    print("machine.%d.log bytes=%d lines=%d" % (node, bytes_written, lines_written))