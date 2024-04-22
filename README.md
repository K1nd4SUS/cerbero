# cerbero

> A packet filtering tool for A/D CTFs.

![GitHub_banner](https://user-images.githubusercontent.com/23193188/142950155-5d2e00a6-7c9f-42db-9cdf-e28783e66f30.gif)

## Table of contents

<!--toc:start-->
- [cerbero](#cerbero)
  - [Table of contents](#table-of-contents)
  - [Intro](#intro)
    - [Why this tool](#why-this-tool)
    - [Is it against the rules to use this tool](#is-it-against-the-rules-to-use-this-tool)
  - [Structure of the project](#structure-of-the-project)
  - [Usage](#usage)
  - [Contributing](#contributing)
<!--toc:end-->

## Intro

This tool is able to filter packets based on their payload by using regular expressions.

### Why this tool

During an A/D, we often had to drop some malicious packets, but to do it properly we had to understand how the service worked and in which programming language it was written. This process is a waste of time. With this tool we can instead drop the malicious packets before they are received from the vulnerable service, making the process simple and implementation agnostic (basically it works like a WAF that process all packets).

### Is it against the rules to use this tool

Probably not, but we do not take any responsibility for its use.

## Structure of the project

This hyper-professional diagram represents on a conceptual level how this tool is structured:

![diagram](https://github.com/K1nd4SUS/cerbero/assets/78105813/35df3ed1-a460-4fa3-ba79-fe169b622477)

## Usage

1. [Deploy cerbero-web with docker compose](/web/README.md#deployment-with-docker-compose)
2. [Download the cerbero binary on the vuln-box](/firewall/README.md#download-the-binary-on-the-vuln-box)

**Warning: cerbero-web must be set up before trying to connect the firewall**, this means that *before starting the `cerbero` binary* you MUST complete the services setup on cerbero-web.

## Contributing

If you wish to contribute to the project, make sure to read the [contributing](/CONTRIBUTING.md) guide first.

