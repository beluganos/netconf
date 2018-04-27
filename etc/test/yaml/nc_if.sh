#! /bin/bash

read -s -p "enter password:" input
echo""

beluganos-interfaces $* | ncclient -P $input edit-config candidate -c -
