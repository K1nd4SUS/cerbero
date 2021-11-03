package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	nfqueue "github.com/florianl/go-nfqueue"
)

func checkFlag(mode string, nfqCoonfig uint16){
	//check if nfqCoonfig is in the allowed range
	if(nfqCoonfig < 1 || nfqCoonfig > 65535){
		fmt.Println("Invalid argument for flag -nfq, the value need to be between 1 and 65535")
		os.Exit(127)
	}

	//check if mode is allowed (must be "w" or "b")
	if(mode != "w" && mode != "b"){
		fmt.Println("Invalid argument for flag -m, must be set to 'w' or 'b'")
		os.Exit(127)
	}
}

func main() {
	// Send ingoing packets to nfqueue queue 100
	// # sudo iptables -I INPUT -p tcp --dport 12345 -j NFQUEUE --queue-num 100

	//flag specifications
	var nfqFlag = flag.Int("nfq", 100, "Queue number")
	var modeFlag = flag.String("m", "w", "Whitelis or Blacklist")
	flag.Parse()

	//some change for the flags
	nfqCoonfig := uint16(*nfqFlag)
	mode := *modeFlag

	//checks flags
	checkFlag(mode, nfqCoonfig)

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
	reg, _ := regexp.Compile(`(ciao)|(s[0-9]+)`)

	fn := func(a nfqueue.Attribute) int {
		id := *a.PacketID
		payload := make([]byte, len(*a.Payload))
		copy(payload, *a.Payload)
		payloadString := string(payload)
		
		
		if(mode == "b"){ //blacklist (if there is a match with the regex, drop the packet)
			if(reg.MatchString(payloadString)){
				log.Print("DROP")
				nf.SetVerdict(id, nfqueue.NfDrop)
			} else {
				log.Print("ACCEPT")
				nf.SetVerdict(id, nfqueue.NfAccept)
			}
			fmt.Printf("%x\n", payloadString)
		}else if (mode == "w"){ //whitelist (if there is a match with the regex, accept the packet)
			if(reg.MatchString(payloadString)){
				log.Print("ACCEPT")
				nf.SetVerdict(id, nfqueue.NfAccept)
			} else {
				log.Print("DROP")
				nf.SetVerdict(id, nfqueue.NfDrop)
			}
			fmt.Printf("%x\n", payloadString)
		}
		return 0
	}

	r := func(e error) int {
		fmt.Println("Error")
		return 0
	}

	err = nf.RegisterWithErrorFunc(ctx, fn, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Block till the context expires
	<-ctx.Done()
}

