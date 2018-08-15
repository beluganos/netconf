# NETCONF/YANG module

**NETCONF/YANG server** to configure following software:

- [Beluganos](https://github.com/beluganos/beluganos)'s main module
- [FRRouting](https://github.com/FRRouting/frr)
- [GoBGP](https://github.com/osrg/gobgp)

The main components of this software are based on [sysrepo](https://github.com/sysrepo/sysrepo) and [Netopeer2](https://github.com/CESNET/Netopeer2). The main efforts of this software is sysrepo client library for FRR, GoBGP, and Beluganos.

## Getting Started

***Important Notice: If you want to just configure Beluganos by NETCONF, you don't have to check this documents.*** This is because this repositories (beluganos/netconf) software will be installed automatically by install tool at [https://github.com/beluganos/beluganos/create.sh](https://github.com/beluganos/beluganos).

If you want to use this software to configure FRRouting or GoBGP, following steps are required.

1. Install software: [install guide](doc/install-guide.md).
2. Setup some configuration: [setup guide](doc/setup-guide.md).
3. Start this program: [operation guide](doc/operation-guide.md).
4. Issue NETCONF commands: [example](doc/examples).

## Support
Github issue page and e-mail are available. If you prefer to use e-mail, please contact `msf-contact-ml [at] hco.ntt.co.jp`.

## Development & Contribution
The main component is written in Go. Any contribution is encouraged. For more details, please refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Beluganos is licensed under the **Apache 2.0** license. See [LICENSE](LICENSE).

## Project
This project is a part of [Multi-Service Fabric](https://github.com/multi-service-fabric/msf).