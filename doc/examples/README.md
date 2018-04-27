## Sample environments

### initial settings
- For PE1: /etc/lxcinit/vpn_mic
- For PE1 VRF10: /etc/lxcinit/vpn_ric

### environments
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
