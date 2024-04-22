# cerbero-firewall

## Table of contents

<!--toc:start-->
- [cerbero-firewall](#cerbero-firewall)
  - [Table of contents](#table-of-contents)
  - [Users documentation](#users-documentation)
    - [Deploy](#deploy)
      - [Download the binary on the vuln-box](#download-the-binary-on-the-vuln-box)
    - [Usage](#usage)
      - [Run with config file](#run-with-config-file)
      - [Run with socket](#run-with-socket)
      - [Help](#help)
<!--toc:end-->

## Users documentation

### Deploy

#### Download the binary on the vuln-box

You can download the latest resease from [here](https://github.com/K1nd4SUS/cerbero/releases).

### Usage

**Warning: The `cerbero` binary needs `root` privileges to run.**

#### Run with config file

The binary can get the configuration from a `json` file (see [config.json](/firewall/config.json) as an example).

```sh
cerbero file <path_to_config> [...options]
```

*The `file` mode is more suited for development/debugging.*

#### Run with socket

When used in conjunction with cerbero-web, you must run the binary in `socket` mode and specify the correct socket `address:port` pair.

```sh
cerbero socket <socket_address:socket_port> [...options]
```

For example, if in cerbero-web (deployed on the same machine as the binary) you decided to run the TCP socket server on port `1234` (specified in the `.env` file by the `SOCKET_PORT` variable) this command will do the job:

```sh
cerbero socket localhost:1234
```

*The `socket` mode is the preferred way to run cerbero in an A/D CTF.*

#### Help

You can find more in-depth information about the commands and flags in the help menu:

```sh
cerbero -h
```

