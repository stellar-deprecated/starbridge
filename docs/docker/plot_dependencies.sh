#!/bin/bash

FILE=$1
NAME=${FILE%.*}
cd /home/user/output/
# analyze file using 2 threads:
python /home/user/poison-ivy/poisonivy.py /home/user/models/$FILE 2
# make png file:
dot -Tpng -o $NAME-invariants.png $NAME.dot
