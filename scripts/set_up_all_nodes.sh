#!/bin/sh

# usage: ./scripts/setup_all.sh <netid>
NETID=$1
REPO=https://github.com/maxh/cmgrep.git

for i in 01 02 03 04 05 06 07 08 09 10; do
    host=fa26-cs425-72$i.cs.illinois.edu
    ssh $NETID@$host "
        git clone $REPO 2>/dev/null || (cd cmgrep && git pull)
        cd cmgrep
        go build -o cmgrep .
        ./scripts/gen_logs.py
        pkill -x cmgrep
        nohup ./cmgrep serve > server.out 2>&1 < /dev/null &
    " &
done

wait
