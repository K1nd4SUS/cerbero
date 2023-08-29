# DOCUMENTATION

## Usage

 il tag nfq viene di default inserito ad ogni servizio (partendo da 100), se si vuole partire da un altro numero, basta inserire il parametro `nfq` con il valore desiderato. I servizi vengono passati da input tramite JSON.

```bash
go run firewall2.go [--nfq 200] [--path ./config.json]
```
oppure
```bash
go build -o firewall
./firewall [--ndq 200] [--path ./config.json]
```