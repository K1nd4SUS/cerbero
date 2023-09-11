# DOCUMENTATION

## Usage

 il tag nfq viene di default inserito ad ogni servizio (partendo da 100), se si vuole partire da un altro numero, basta inserire il parametro `nfq` con il valore desiderato. I servizi vengono passati da input tramite JSON.

```bash
go run firewall2.go [--nfq 200] [--path ./config.json]
```
oppure
```bash
go build -o firewall
./firewall [--nfq 200] [--path ./config.json]
```

> aggiungere script per test automatico del corretto setup
> evitare in ogni modo il crash (es typo nella config salta il servizio e manda warning)
> Per filtrare docker usare la chain DOCKER-INGRESS, mentre per servizi non docker INPUT