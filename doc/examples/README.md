## Sample configurations

This folder contains the sample xml file of NETCONF's `<edit-config>`.

### Initial settings

- For PE1: `/etc/lxcinit/vpn_mic`
- For PE1 VRF10: `/etc/lxcinit/vpn_ric`

### Items of sample configurations

1. `clear-<module-name>.xml`
	- Init each modules
	- Default-operation: replace

1. `routing-policy.xml`
	- Add GoBGP routing policy
	- Default-operation: replace or merge

1. `interfaces.xml`
	- Add physical interface settings (eth1~5) to interfaces modules
	- Default-operation: replace or merge

1. `network-instance-mic.xml`
	- Add network instance (PE1) without interface, ospf, bgp, and mpls settings
	- Default-operation: replace or merge

1. `interfaces_add_subif.xml`
	- Add sub-IF settings (eth1.100, eth2.100, eth4.20, eth5.20) to physical interfaces
	- Default-operation: merge
   
1. `network-instance-mic_add_ifaces.xml`
	- Add interfaces (eth1.100, eth2.100) to network-instance (PE1), add interface IP address to interfaces (eth1.100, eth2.100)
   - Default-operation: merge

1. `network-instance-mic_add_protos.xml`
	- Add ospf, bgp, and mpls settings
   - Default-operation: merge

1. `network-instance-vrf10.xml`
	- Add new VRF (PE-VRF10) and many settings to network instance
   - Default-operation: merge

1. `network-instance-vrf10_del_iface.xml`
   - Delete interfaces (eth5.20) from network-instance (PE1-VRF10)
   - Default-operation: merge

### Environments of this sample configurations
~~~
                                   - - ---+
                                          |
     AS: 65001         AS: 65001      AS: 65001
     +------+          +------+       +------+
     |   P  |          |   P  |       |  RR  | 10.0.0.100
     +------+          +------+       +------+
         \               /
          \             /
        eth1.100   eth2.100
    10.10.1.6/30   10.10.2.6/24
          +-------------+
          |             | RD: - / RT: 10:1 (dummy)
          |     PE1     | AS: 65001 
          |   10.0.1.6  | import-policy: -
          |             | export-policy: policy-local-pref, policy-next-hop-self
          +-------------+
          +-------------+
          |             | RD: 10:2001 / RT: 10:1
          |  PE1-VRF10  | AS: 65001
          |   10.0.1.6  | import-policy: policy-local-pref-vrf
          |             | export-policy: policy-local-pref
          +-------------+
        eth4.20    eth5.20
   30.10.1.1/24    30.10.2.1/24
          /             \
         /               \
   30.10.1.2/24      30.10.2.2/24
     +------+          +------+
     |  CE  |          |  CE  |
     +------+          +------+
      AS: 30
~~~

### Operations of NETCONF

NETCONF over ssh will utilize TCP 830 port. NETCONF session will be started by following commands.

```
$ ssh -s <server-ip> -p 830 netconf
```

After exchanging `<hello>` message, you can operate `<get-config>` or `<edit-config>` operations. The example log is located at `operation.log`.
