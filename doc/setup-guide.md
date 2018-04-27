# Setup guide

## Pre-requirements
- Install required components
	- Please refer [install-guide.md](install-guide.md).

## Resources

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

## 2. create lxcinit.sh

Once network-instance is created, `lxcinit.sh` will be worked. You can customize   this script.

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
