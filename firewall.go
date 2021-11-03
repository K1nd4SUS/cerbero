package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"regexp"
	"flag"
	nfqueue "github.com/florianl/go-nfqueue"
)

func main() {
	// Send ingoing packets to nfqueue queue 100
	// # sudo iptables -I INPUT -p tcp --dport 12345 -j NFQUEUE --queue-num 100

	//flag specifications
	var nfq = flag.Int("nfq", 100, "Queue number")
	flag.Parse()
	nfqCoonfig := uint16(*nfq)

	// Set configuration options for nfqueue
	config := nfqueue.Config{
		NfQueue:      nfqCoonfig,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 15 * time.Millisecond,
	}

	nf, err := nfqueue.Open(&config)
	if err != nil {
		fmt.Println("could not open nfqueue socket:", err)
		return
	}
	defer nf.Close()

	ctx:= context.Background()
	reg, _ := regexp.Compile("ciao")

	fn := func(a nfqueue.Attribute) int {
		id := *a.PacketID
		payload := make([]byte, len(*a.Payload))
		copy(payload, *a.Payload)
		payloadString := string(payload)
		
		if(reg.MatchString(payloadString)){
			log.Print("DROP")
			nf.SetVerdict(id, nfqueue.NfDrop)
		} else {
			log.Print("ACCEPT")
			nf.SetVerdict(id, nfqueue.NfAccept)
		}
		fmt.Printf("%x\n", payloadString)
		
		return 0
	}

	r := func(e error) int {
		fmt.Println("Error")
		return 0
	}

	// Register your function to listen on nflqueue queue 100
	err = nf.RegisterWithErrorFunc(ctx, fn, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Block till the context expires
	<-ctx.Done()
}

