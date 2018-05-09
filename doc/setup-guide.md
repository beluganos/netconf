# Setup guide

This document shows how to prepare NETCONF server to recieve the NETCONF requests.

## Pre-requirements
- You have to finish installing before setup. Please refer [install-guide.md](install-guide.md) before proceeding.

## Resources

The strings like `std_mic` and `vpn_ric` represents the type of network-instance.

~~~
    etc/
        lxcinit/
            std_mic/
            std_ric/
            vpn_mic/
            vpn_ric/
                lxcinit.sh
                conf/
                service/
~~~


## Understand operating principle

This repositories software has been developed based on the idea of [OpenConfig](http://www.openconfig.net/). This is because 
at least one "**network-instance**" is required by this software. For example, in general IP routers, single **network-instance** (type = `DEFAULT_INSTANCE`) exists is assumed.

In this section, the operation principle about adding network-instance is described. Generally, **adding network-instance may be operated at first NETCONF request**, except in case of adding beforehand some network-instance by another methods.

### (step 0) NETCONF request

For example, following edit-config requests (default operation is "merge") is assumed.

```
<network-instances xmlns="https://github.com/beluganos/beluganos/yang/network-instance">
  <network-instance>
    <name>master-instance</name>
    <config>
      <name>master-instance</name>
      <type xmlns:oc-ni-types="http://openconfig.net/yang/network-instance-types">oc-ni-types:DEFAULT_INSTANCE</type>
      <router-id>10.0.1.6</router-id>
    </config>
    <loopbacks>
      <loopback>
        <id>lo</id>
        <config>
          <id>lo</id>
        </config>
        <addresses>
          <address>
            <index>0</index>
            <config>
              <index>0</index>
              <ip>10.0.1.6</ip>
              <prefix-length>32</prefix-length>
            </config>
          </address>
        </addresses>
      </loopback>
    </loopbacks>
    <interfaces/>
    <protocols/>
  </network-instance>
</network-instances>
```

### (step 1) execute initialization script at host

The script of `lxcinit.sh` is executed automatically at host server (**NOT** linux container). The argument is container's name, path of directory, and the strings of "local".

```
(host)$ ~/etc/lxcinit/std_mic/lxcinit.sh master-instance ~/etc/lxcinit/std_mic local
```

This operation is owned by nc-module.

### (step 2) transfer files

Required files which exists on `etc/lxcinit/<container-type>` is tranfered to `/tmp/<container-type>` at linux container.

### (step 3) execute initialization script at linux container

The argument is container's name, path of directory.

```
(container)$ /tmp/lxcinit.sh master-instance /tmp/std_mic
```

---

## 1. Specify instance type

You have to specify instance type to create network instance, nevertheless you will use just only single network instance. To specify, the field of `type` and `route-target` are used.

```
module: beluganos-network-instance
    +--rw network-instances
       +--rw network-instance* [name]
          +--rw name          -> ../config/name
          +--rw config
          |  +--rw type?                  identityref (*)
          |  +--rw route-target?          oc-ni-types:route-distinguisher (*)
          |  +--rw ....
```

The supported network instance (i.e. Linux container) type is following: 

| Type             | Route-target | Name    | Description        |
| ---------------- | ------------------------------------ | ---- | ---- |
| DEFAULT_INSTANCE | No | std_mic | Standard network instance            |
| L3VRF            | No | std_ric | Virtual router (VRF-Lite)            |
| DEFAULT_INSTANCE | Yes(*1)| vpn_mic | Standard network instance with L3VPN  |
| L3VRF            | Yes| vpn_ric | VRF for L3VPN                        |

(*1) As for any value, please fill it. This value is not used by Beluganos.

## 2. Edit initialization script

Once network-instance is created, `lxcinit.sh` will be worked. You can customize this script. For more detail, please refer the section of "Understand operating principle" in this document.

## 3. Setup systemd

If you want use systemd service, please put following files at `/etc/lxcinit/`.

```
/etc/lxcinit/
        std_mic/
        std_ric/
        vpn_mic/
        vpn_ric/
            lxcinit.sh
            conf/
            service/
```

## Next steps

Please refer [operation-guide.md](operation-guide.md).