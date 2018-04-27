#! /bin/bash

read -s -p "enter password:" input
echo""

beluganos-routing-policy $* | ncclient -P $input edit-config candidate -c -
