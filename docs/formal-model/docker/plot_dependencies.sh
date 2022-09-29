#!/bin/bash

trap "exit" INT

FILE_PATH=$1
FILENAME=${FILE_PATH##*/}
NAME=${FILENAME%.*}
# MODEL_DIR=$(dirname $FILE_PATH)
OUTPUT_DIR=$2
N_THREADS=$3
cd $OUTPUT_DIR
python /home/user/poison-ivy/poisonivy.py $FILE_PATH $N_THREADS &
wait
# make png file:
dot -Tpng -o $NAME-invariants.png $NAME.dot
echo "finished"
