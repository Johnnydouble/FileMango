#!/bin/bash

#run program and get pid
/home/test/FMbundle/FileMango >> ./output.txt &
FM_PID=$!

#log stuff
while sleep 1; do ps --no-headers -o '%cpu,%mem' -p $FM_PID >> ./log.txt; done &
while sleep 1; do date "+%T" >> ./output.txt && date "+%T" >> ./log.txt; done &
#wait to do final tests
sleep 100

#do final tests
echo "file 1" >> ./output.txt
echo "file 1" >> ./log.txt
wget https://upload.wikimedia.org/wikipedia/commons/a/af/Tux.png -P ~/Downloads
sleep 5

echo "file 2" >> ./output.txt
echo "file 2" >> ./log.txt
wget https://live.staticflickr.com/3903/15218475961_963a4c116e_b.jpg -P ~/Downloads
sleep 5

#done
echo "done" >> ./output.txt
echo "done" >> ./log.txt
echo "done"
