# Cerbero - A packet filtering tool

![GitHub_banner](https://user-images.githubusercontent.com/23193188/142950155-5d2e00a6-7c9f-42db-9cdf-e28783e66f30.gif)

---
### Description
This tool can filter packets based on the payload. The match is made with regex with two modalities, whitelist and blacklist.

### Why this tool?
During an A/D, we often had to drop some malicious packets, but to do it properly we had to understand how the service worked and in which programming language it was written. This process is a waste of time. With this tool we can instead drop the malicious packets before they are received from the vulnerable service, making the process simple and implementation agnostic (basically it works like a WAP that process all packets).

### Why GO?
Because [MΛRC02U1F4A4](https://github.com/M4RC02U1F4A4) wanted to learn GO, and also the library [nfqueue](https://pkg.go.dev/github.com/florianl/go-nfqueue) was easy to use.

### Is it against the rules to use this tool? 
Probably not, but we do not take any responsibility for its use .

---

### How to use the tool

`⚠️ The program need root permission to run ⚠️`

1. Download the last resease from [here](https://github.com/K1nd4SUS/Cerbero/releases)
2. Choose between two starting method: CLI or JSON

```
Usage of ./cerbero:
  -dport int
        Destination port number (default 8080)
  -mode string
        Whitelist(w) or Blacklist(b) (default "b")
  -nfq int
        Queue number (default 100)
  -p string
        Protocol 'tcp' or 'udp' (default "tcp")
  -path string
        Path to the json config file (default "./config.json")
  -r string
        Regex to match, follow this format: '(regex1)|(regex2)|...'
  -t string
        Type of input, 'j' for json or 'c' for command line (if j is choosen, only the path flag is considered (default "c")
```

### CLI
Example
```console
$ sudo ./cerbero -t c -dport 12345 -mode b -nfq 101 -p udp -r '(malicious)'
```
In this way the tools is going to filter all the packets with destination port `12345`, protocol `udp` and that contain the word `malicious`

### JSON
```console
$ sudo ./cerbero -t j
```
Sample configuration file 
```json
{
    "services": [
      {
        "name": "Service 1",
        "nfq": 100,
        "mode": "b",
        "protocol": "udp",
        "dport": 8080,
        "regexList": [
          "(malicius)"
        ]
      },
      {
        "name": "Service 2",
        "nfq": 101,
        "mode": "w",
        "protocol": "udp",
        "dport": 8181,
        "regexList": [
          "(safe)"
        ]
      }
    ]
}
```
---

### A/D Setup - DOCKER

1. Clone the repo
```console
local$ git clone git@github.com:K1nd4SUS/Cerbero.git
```
2. Change folder
```console
local$ cd Cerbero
```
3. Compile `cerbero`
```console
local$ go build -o cerbero firewall.go
```
4. Edit `config.json` with the services informations
5. Insert into `docker-compose.yml` the IP (`VULNBOX_IP`)
6. Generate a pair of ssh keys (`ssh-keygen -t ed25519`) to allow the container to authenticate with the vulnbox
7. Start the container
```console
local$ docker-compose up --build -d
```
7. All files will have appeared on the vulnbox
```console
vulnbox$ ls
cerbero
vulnbox$ cd cerbero/
vulnbox$ ls
cerbero  config.json  start_cerbero
```
8. Start `cerbero` on the `vulnbox`
```console
vulnbox$ ./start_cerbero
```
9. Access the WebGUI on your local machine via `localhost:51645`

### Point 7 is very important, DON'T USE THE KEYS IN THE REPO!!!

---

> Developed by [MΛRC02U1F4A4](https://github.com/M4RC02U1F4A4) \
> Thanks to [Tiziano Radicchi](https://github.com/tiz314) for the help with the web interface
