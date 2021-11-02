package main

import (
	"context"
	"fmt"
	"time"

	nfqueue "github.com/florianl/go-nfqueue"
)

func main() {
	// Send outgoing pings to nfqueue queue 100
	// # sudo iptables -I OUTPUT -p tcp --dport 12345 -j NFQUEUE --queue-num 100

	// Set configuration options for nfqueue
	config := nfqueue.Config{
		NfQueue:      100,
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fn := func(a nfqueue.Attribute) int {
		// Just print out the id and payload of the nfqueue packet
		payload := *a.Payload
		fmt.Printf("%v\n", payload)
		// if(a.Payload){
		// 	nf.SetVerdict(id, nfqueue.NfAccept)
		// } 
		
		return 0
	}

	// Register your function to listen on nflqueue queue 100
	err = nf.Register(ctx, fn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Block till the context expires
	<-ctx.Done()
}