# DOCUMENTATION

## Usage
Output of `./firewall2 --help`:
```
Usage of ./firewall2:
  -docker string
    	select iptables chain list (default "INPUT")
  -nfq int
    	Queue number (optional, default 100 onwards) (default 100)
  -path string
    	Path to the json config file (default "./config.json")
```

The tag "nfq" is by default set to 100 onwards (so if you start Cerbero with 10 services, 100 through 109 will be used). If you want to start another Cerbero instance, you have to change this number or else it will conflict with the previous iptables rules. For instance, run:
```bash
go run firewall2.go [--nfq 200] [--path ./config.json]
```
or
```bash
go build ./firewall2.go
./firewall2 [--nfq 200] [--path ./config.json]
```

> aggiungere script per test automatico del corretto setup
> evitare in ogni modo il crash (es typo nella config salta il servizio e manda warning)
> Per filtrare docker usare la chain DOCKER-INGRESS, mentre per servizi non docker INPUT
> aggiungere: I servizi vengono passati da input tramite JSON.