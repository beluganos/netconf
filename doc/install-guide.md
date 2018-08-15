# Install guide

## Pre-requirements

- OS:  Ubuntu 18.04 server edition
- Mem: 4GB+
- HDD: 10GB+ (each LXC needs extra 1GB storage)
- Software: [sysrepo(v0.7.1)](https://github.com/sysrepo/sysrepo/releases/tag/v0.7.1), [Netopeer2(v0.4-r1)](https://github.com/CESNET/Netopeer2/releases/tag/v0.4-r1)

## 1. Install beluganos-netconf

```
$ git clone https://github.com/beluganos/netconf && cd netconf
```

### Set internet proxy

**For proxy environments only:** if you need use proxy server to connect internet, please comment out and edit PROXY settings in `create.ini`.

```
$ vi create.ini

PROXY=http://192.168.1.100:8080
```

### Change netopeer2/sysrepo build options

Changing building option of netopeer2/sysrepo is recommended. If you did not change this settings, you cannot use systemctl commands to start netopeer2/sysrepo.

```
$ vi etc/sysrepo/sr_install.sh

CMAKE_BUILD="-DCMAKE_BUILD_TYPE=Debug"
# CMAKE_BUILD="-DCMAKE_BUILD_TYPE=Release"
```

### Build

```
$ ./create.sh
```

## 2.  Register as a service

```
$ sudo make install-service
```

---

## 3. Install LXD containers

If you want to configure Linux's base settings by NETCONF, following steps are required. 

In this software, "network-instance" means "linux container (LXD)". This means that if you execute the request to create a "network-instance", a "linux container (LXD)" will be created. This new LXD instance is a just copy of "base" image. Thus, you should create "base" image before issue NETCONF commands.

### Current support status

- OS: Ubuntu 18.04
- ARCH: x86_64

### Create "base" image

```
$ lxc launch ubuntu:18.04 temp
$ lxc stop temp
$ lxc publish temp --alias base
$ lxc delete temp
$ lxc image list
+--------------+--------------+--------+---------------+--------+----------+---------------+
|    ALIAS     | FINGERPRINT  | PUBLIC |  DESCRIPTION  |  ARCH  |   SIZE   |  UPLOAD DATE  |
+--------------+--------------+--------+---------------+--------+----------+---------------+
| base         | 01234567890  | no     |  foo bar      | x86_64 | ---.--MB | -- / -- / --  |
+--------------+--------------+--------+---------------+--------+----------+---------------+
```

### GoBGP/FRR

If you need to configure GoBGP/FRR, you should setup this software before starting netopeer2/sysrepo. Please note that when network-instance is added/deleted by NETCONF, GoBGP/FRR settings will be changed. 

You have two options for installation:

1. Installation to "base" container in advance
	- The network-instance will be created based one "base" container.
1. Use `lxcinit.sh`
	- Once network-instance is created, `lxcinit.sh` will be worked. You can customize this script. For more detail, please refer [setup-guide.md](setup-guide.md).

## Next steps

Please refer [setup-guide.md](setup-guide.md).