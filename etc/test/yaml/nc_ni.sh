#! /bin/bash

read -s -p "enter password:" input
echo""

beluganos-network-instance $* | ncclient -P $input edit-config candidate -c -
