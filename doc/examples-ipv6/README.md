## Sample configurations

This folder contains the sample xml file of NETCONF's `<edit-config>`.

### Initial settings

- For R1: `/etc/lxcinit/std_mic`

### Items of sample configurations

1. `clear-<module-name>.xml`
        - Init each modules
        - Default-operation: replace

1. `routing-policy.xml`
        - Add GoBGP routing policy
        - Default-operation: replace or merge

1. `interfaces.xml`
        - Add physical interface settings (eth1~3) to interfaces modules
        - Default-operation: replace or merge

1. `network-instance.xml`
        - Add network instance (R1) without interface, ospf, and bgp settings
        - Default-operation: replace or merge

1. `interfaces_add_subif.xml`
        - Add sub-IF settings (eth1, eth2, eth3) to physical interfaces
        - Default-operation: merge
   
1. `network-instance_add_ifaces.xml`
        - Add interfaces (eth1, eth2, eth3) to network-instance (R1), add interface IP address to interfaces(eth1, rth2, eth3).
   - Default-operation: merge

1. `network-instance_add_protos.xml`
        - Add ospf, ospfv3 and bgp settings
   - Default-operation: merge

### Environments of this sample configurations
~~~

     
               +----------+
               |    R2    | ospf
               | 10.0.0.2 | ospfv3
               |          | Lo(6): 2001::2/128
               +----------+
                    | .1
                    | 10.2.0.0/24
                    | 2001:db8:2::/64
                    | .2
                   eth1
             +-------------+
             | <beluganos> | AS: 65001
             |      R1     | import-policy: -
             |   10.0.0.1  | export-policy: policy-next-hop-self
             |             | Lo(6): 2001::1/128
             +-------------+
             eth2       eth3
             / .1          \ .1
    2001:db8:3::/64     10.4.0.0/24
          / .2                \ .2
    +----------+         +----------+
    |    R3    |         |    R4    |
    | 10.0.0.3 |         | 10.0.0.4 |
    +----------+         +----------+
     AS: 65001            AS: 65001
     Lo(6): 2001::3/128

~~~

### Operations of NETCONF

NETCONF over ssh will utilize TCP 830 port. NETCONF session will be started by following commands.

```
$ ssh -s <server-ip> -p 830 netconf
```

After exchanging `<hello>` message, you can operate `<get-config>` or `<edit-config>` operations. The example log is located at `operation.log`.
