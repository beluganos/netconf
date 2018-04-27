# beluganos-netconf tool

## 1.  How to install

```
> make install-python
```

## 2. How to use.

### 2.1. generate XML (beluganos-interfaces module)

```
> beluganos-interfaces <yaml file>
```

#### 2.1.1. yaml format

```
---

network:
  ethernets:
    eth1: {}
    eth2:
      mac-address: "00:22:33:44:55:66"
    eth3:
      mac-address: "00:77:88:99:aa:bb"
      addresses:
        - 30.10.1.6/30

  vlans:
    eth1.100:
      id: 100
      link: eth1
      mtu: 1500
      addresses:
        - 10.10.1.6/30

    eth2.100:
      id: 100
      link: eth2
      mtu: 1500
      addresses:
        - 20.10.1.6/30
```

### 2.2. generate XML (beluganos-routing-policy module)

```
> beluganos-routing-policy <yaml file>
```

### 2.2.1. yaml format

*Only set-next-hop and set-local-pref actions are supported.*

```
---

policy-definitions:
  policy-next-hop-self:       # policy-name
    stmts:
      stmt-next-hop-self:
        actions:
          policy-result: ACCEPT_ROUTE
          bgp:
            set-next-hop: SELF
            
  policy-local-pref:
    stmts:
      stmt-local-pref:
        actions:
          policy-result: ACCEPT_ROUTE
          bgp:
            set-local-pref: 100
```

### 2.3. generate XML (beluganos-network-instance module)

```
> beluganos-network -instance <yaml file>
```

### 2.3.1. yaml format

```
---

network-instances:
  PE1:                            # container name
    type: DEFAULT_INSTANCE        # DEFAULT_INSTANCE or L3VRF
    router-id: 10.0.1.6
    route-distinguisher: "10:1"
    route-target: "10:2001"

    loopbacks:
      lo:
        - 10.0.1.6/32

    interfaces:
      eth1: []
      eth1.100: [both] # empty, iface, subif, both
      eth2: []
      eth2.100: [iface, subif]

    routes:
      30.10.2.0/24:
        - via: 30.10.2.2
      30.10.3.0/24:
        - dev: eth5.10
      30.10.4.0/24:
        - drop: true

    ospfv2:
      router-id: 10.0.1.6
      areas:
        0.0.0.0:              # area-id
          lo:
            metric: 10
            passive: true
          eth1.100: {}
          eth2.100:
            metric: 200
            passive: false
            hello-interval: 40
            dead-interval: 10

    mpls:
      ldp:
        router-id: 10.0.1.6
        null-label: EXPLICIT  # EXPLICIT, IMPLICIT, None
        ipv4:
          transport-address: 10.0.1.6
          session-holdtime: 3600
          egress-explicit-null: true
        timers:
          hello-holdtime: 180
          hello-interval: 10
        interfaces:
          eth1.100: {}
          eth2.100:
            hello-holdtime: 180
            hello-interval: 10

    bgp:
      router-id: 10.0.1.6
      as: 65001
      zebra:
        enabled: false              # true or false
      neighbors:
        10.0.0.100:
          peer-as: 65001
          # local-as: 65001
          transport-local-address: 10.0.1.6
          afi-safis:
            - L3VPN_IPV4_UNICAST    # IPV4_UNICAST or L3VPN_IPV4_UNICAST
          timers:
            hold-time: 3600
            keepalive-interval: 180
          apply-policy:
            import-policy: []
            export-policy:
              - policy-next-hop-self
              - policy-local-pref
            default-import-policy: ACCEPT_ROUTE
            default-export-policy: ACCEPT_ROUTE

```
