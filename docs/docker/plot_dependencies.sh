#!/bin/bash

FILE=$1
NUM_PROC=$2
NAME=${FILE%.*}
cd /home/user/output/
# analyze file using NUM_PROC threads:
python /home/user/poison-ivy/poisonivy.py /home/user/models/$FILE $NUM_PROC
# make png file:
dot -Tpng -o $NAME-invariants.png $NAME.dot
