# Operation guide

The start/stop operation is described at this documents.

## Pre-requirements

- You have to finish installation and setup. Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.


## Start / Stop

You have two options to operate this software.

### (Option A) Launch as a service

```
$ sudo systemctl start netopeer2-server
$ sudo systemctl start ncm.target

$ sudo systemctl stop ncmd
$ sudo systemctl stop ncmi
$ sudo systemctl stop ncms
```

### (Option B) CLI

You need two terminals, because this launching scripts will take your standard output. Generally, this option is used for debugging is assumed.

```
$ cd ~/netconf/etc/test
$ ./nop2.sh    # terminal-1
$ ./ncmall.sh  # terminal-2
```