# Cerbero - A packet filtering tool

![GitHub_banner](https://user-images.githubusercontent.com/23193188/140656404-5ae9dc4a-ac35-40f0-aa7d-b2869e694a28.png)

---
### Description
This tool can filter packets based on the payload. The match is made with regex with two modalities, whitelist and blacklist.

### Why this tool?
During an A/D, we often had to drop some malicious packets, but to do it properly we had to understand how the service worked and in which programming language it was written. This process is a waste of time. With this tool we can instead drop the malicious packets before they are received from the vulnerable service, making the process simple and implementation agnostic (basically it works like a WAP that process all packets).

### Why GO?
Because [M4RC02U1F4A4](https://github.com/M4RC02U1F4A4) wanted to learn GO, and also the library [nfqueue](https://pkg.go.dev/github.com/florianl/go-nfqueue) was easy to use.

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
$ sudo ./cerbero -t c -dport 12345 -mode b -nfq 101 -p udp -r '(malicius)'
```
In this way the tools is going to filter all the packets with destination port `12345`, protocol `udp` and that contain the word `maliciust`

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
