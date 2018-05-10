# Operation guide

The start/stop operation is described at this documents.

## Pre-requirements

- You have to finish installation and setup. Please refer [install-guide.md](install-guide.md) and [setup-guide.md](setup-guide.md) before proceeding.


## Start / Stop

You have two options to operate this software. Generally, launching as a linux service is recommended.

### (Option A) Launch as a service

To start, 

```
$ sudo systemctl start netopeer2-server
$ sudo systemctl start ncm.target
```

To stop,

```
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

## NETCONF operation

In default settings, SSH port of NETCONF is 830. For more detail, please refer [examples/README.md](examples/README.md).