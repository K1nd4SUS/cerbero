# Cerbero - A packet filtering tool

![GitHub_banner](https://user-images.githubusercontent.com/23193188/142950155-5d2e00a6-7c9f-42db-9cdf-e28783e66f30.gif)

---
### Description
This tool is able to filter packets based on the payload by using regex.

### Why this tool?
During an A/D, we often had to drop some malicious packets, but to do it properly we had to understand how the service worked and in which programming language it was written. This process is a waste of time. With this tool we can instead drop the malicious packets before they are received from the vulnerable service, making the process simple and implementation agnostic (basically it works like a WAP that process all packets).

### Is it against the rules to use this tool? 
Probably not, but we do not take any responsibility for its use.

---

### How to use the tool

`⚠️ The program need root permission to run ⚠️`

You can download the last resease from [here](https://github.com/K1nd4SUS/Cerbero/releases)

```
Usage of ./cerbero:
  -chain string
        Input chain name. (default "INPUT")
  -config string
        Relative or absolute path to the JSON configuration file. (default "./config.json")
  -v    Enable DEBUG-level logging.
```
