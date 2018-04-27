# Install guide

## Pre-requirements

- OS:  Ubuntu 17.10 Server edition
- Mem: 4GB+
- HDD: 10GB+ (each LXC needs extra 1GB HDD)
- Software: 

## 1. Install beluganos-netconf

```
$ git clone https://github.com/beluganos/netconf
$ cd netconf
$ ./create.sh
```

### 1.1. Configure internet proxy

**For proxy environments only:** if you need use proxy server to connect internet, please comment out and edit PROXY settings in `create.ini`.

```
$ vi create.ini

PROXY=http://192.168.1.100:8080
```

### 1.2. Change building settings of Netopeer2 / Sysrepo

```
$ vi etc/sysrepo/sr_install.sh

CMAKE_BUILD="-DCMAKE_BUILD_TYPE=Debug"
# CMAKE_BUILD="-DCMAKE_BUILD_TYPE=Release"
```

## 2. Configure systemctl

```
$ sudo make install-services
```

---

## 3. Install LXD containers

If you want to configure Linux's base settings by NETCONF, following steps are required. 

In this software, "network-instance" means "linux container (LXD)". This means that if you issue the request to create a "network-instance", a "linux container (LXD)" will be created. This new LXD instance is a just copy of "base" image. Thus, you should create "base" image before issue NETCONF commands.

### Current support status

- OS: Ubuntu 17.10
- ARCH: x86_64

### Create "base" image

```
$ lxc launch ubuntu:17.10 temp
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