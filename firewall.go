package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"time"

	nfqueue "github.com/florianl/go-nfqueue"
)

func checkFlag(mode string, nfqCoonfig uint16, protocol string, port int, inType string, path string){
	//check if nfqCoonfig is in the allowed range
	if(nfqCoonfig < 1 || nfqCoonfig > 65535){
		fmt.Println("Invalid argument for flag -nfq, the value need to be between 1 and 65535")
		os.Exit(127)
	}

	//check if mode is allowed (must be "w" or "b")
	if(mode != "w" && mode != "b"){
		fmt.Println("Invalid argument for flag -mode, must be set to 'w' or 'b'")
		os.Exit(127)
	}

	//checks if the procols is correct (must be "tcp" or "udp")
	if(protocol != "tcp" && protocol != "udp"){
		fmt.Println("Invalid argument for flag -p, must be set to 'tcp' or 'udp'")
		os.Exit(127)
	}

	//check if the port number is right
	if(port < 1 || port > 65535){
		fmt.Println("Invalid argument for flag -dport, the value need to be between 1 and 65535")
		os.Exit(127)
	}

	//check if the input type is righe
	if(inType != "c" && inType != "j"){
		fmt.Println("Invalid argument for flag -t, must be set to 'c' or 'j'")
		os.Exit(127)
	}

	_, err := os.Open(path)
	if (err != nil){
		fmt.Println("File not found")
		os.Exit(127)
	}
}

func execJson(path string){
	//loop for create iptables rules
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		<-c
		fmt.Println("\nRemoving iptables rule")
		//loop for delete iptables rules
		//cmd := exec.Command("iptables", "-D", "INPUT", "-p", protocol, "--dport", strconv.FormatInt(int64(port), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(nfqCoonfig), 10))
		//cmd.Run()
		fmt.Println("Done!")
		os.Exit(0)
	}()	
	var wg sync.WaitGroup
	wg.Add(2)
	//loop for start the go routines with exeC
	wg.Wait()

}

func exeC(mode string, nfqCoonfig uint16, protocol string, port int){

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

func main() {
	// Send ingoing packets to nfqueue queue 100
	// $ sudo iptables -I INPUT -p tcp --dport 12345 -j NFQUEUE --queue-num 100

	//flag specifications
	var inTypeFlag = flag.String("t", "c", "Type of input, 'j' for json or 'c' for command line (if j is choosen, only the path flag is considered")
	var pathFlag = flag.String("path", "./config.json", "Path to the json config file")
	var nfqFlag = flag.Int("nfq", 100, "Queue number")
	var modeFlag = flag.String("mode", "b", "Whitelist(w) or Blacklist(b)")
	var protocolFlag = flag.String("p", "tcp", "Protocol 'tcp' or 'udp'")
	var dportFlag = flag.Int("dport", 8080, "Destination port number")
	flag.Parse()

	//some change for the flags
	inType := *inTypeFlag
	nfqCoonfig := uint16(*nfqFlag)
	mode := *modeFlag
	protocol := *protocolFlag
	port := *dportFlag
	path := *pathFlag

	//checks flags
	checkFlag(mode, nfqCoonfig, protocol, port, inType, path)

	if(inType == "j"){
		execJson(path)
	} else {
		//capture ctrl+c
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func(){
			<-c
			fmt.Println("\nRemoving iptables rule")
			cmd := exec.Command("iptables", "-D", "INPUT", "-p", protocol, "--dport", strconv.FormatInt(int64(port), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(nfqCoonfig), 10))
			cmd.Run()
			fmt.Println("Done!")
			os.Exit(0)
		}()	
		//ADD iptables rule
		cmd := exec.Command("iptables", "-I", "INPUT", "-p", protocol, "--dport", strconv.FormatInt(int64(port), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(nfqCoonfig), 10))
		_, err := cmd.Output()
		if err != nil {
			fmt.Println("The program must be run as root")
			os.Exit(126)
		}
		exeC(mode, nfqCoonfig, protocol, port)
	}

}

