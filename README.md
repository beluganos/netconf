# NETCONF/YANG module

**NETCONF/YANG server** to configure following software:

- [Beluganos](https://github.com/beluganos/beluganos)'s main module
- [FRRouting](https://github.com/FRRouting/frr)
- [GoBGP](https://github.com/osrg/gobgp)

The main components of this software are based on [sysrepo](https://github.com/sysrepo/sysrepo) and [Netopeer2](https://github.com/CESNET/Netopeer2). The main efforts of this software is sysrepo client library for FRR, GoBGP, and Beluganos.

## Getting Started

***If you want to use this software in order to just configure Beluganos, you don't have to check this repositories.*** This is because this repositories software will be installed automatically by install tool at [https://github.com/beluganos/beluganos](https://github.com/beluganos/beluganos).

If you want to use this software to configure not Beluganos but FRRouting or GoBGP, following steps are required.

1. Install software: [install guide](doc/install-guide.md).
2. Setup some configuration: [setup guide](doc/setup-guide.md).
3. Start this program: [operation guide](doc/operation-guide.md).
4. Issue NETCONF commands: [example](doc/examples).

## Support
Github issue page and e-mail are available. If you prefer to use e-mail, please contact `msf-contact-ml [at] hco.ntt.co.jp`.

## Development & Contribution
Any contribution is encouraged. The main component is written in Go. If you wish to create pull-request on github.com, please kindly create your request for **develop branch**, not master branch. If you find any issue, please kindly notify us by github issue pages.

For more details, please refer to [CONTRIBUTING.md](CONTRIBUTING.md).

## License
Beluganos is licensed under the **Apache 2.0** license. See [LICENSE](LICENSE).

## Project
This project is a part of [Multi-Service Fabric](https://github.com/multi-service-fabric/msf).