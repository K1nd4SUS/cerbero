package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	nfqueue "github.com/florianl/go-nfqueue"
)

type Services struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name		string 		`json:"name"`
	Nfq 		int 		`json:"nfq"`
	Mode 		string 		`json:"mode"`
	Protocol	string 		`json:"protocol"`
	Dport 		int 		`json:"dport"`
	RegexList	[]string	`json:"regexList"`
}

//check params validity
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

	if (inType == "j"){
		_, err := os.Open(path)
		if (err != nil){
			fmt.Println("File not found")
			os.Exit(127)
	}
	}
	
}

//hash a string, given file path
func hash(path string) (hash string){
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

//onModify
func watchFile(path string, canale chan string){
	oldHash := hash(path)
	for(true){
	time.Sleep(5 * time.Second)
	newHash := hash(path)
	if (oldHash != newHash){
		canale <- "-"
		fmt.Println("File edited")
	}
	oldHash = newHash
	}
}

//serialyze input
func readJson(path string)(Services){
	jsonFile, _ := os.Open(path)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var services Services
	json.Unmarshal(byteValue, &services)
	return services
}

func execJson(path string){
	services := readJson(path)
	//loop for create iptables rules
	for k:= 0; k<len(services.Services); k++{
		cmd := exec.Command("iptables", "-I", "INPUT", "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
		cmd.Run()
	}
	//prepare oninterrupt event
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		<-c
		fmt.Println("\nRemoving iptables rule")
		//loop for delete iptables rules
		for k:= 0; k<len(services.Services); k++{
			cmd := exec.Command("iptables", "-D", "INPUT", "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
			cmd.Run()
		}
		fmt.Println("Done!")
		os.Exit(0)
	}()	
	//start waitgroup 
	var wg sync.WaitGroup
	//onmodify for json
	alertFileEdited := make(chan string)
	wg.Add(len(services.Services)+1)
	//loop for start the go routines with exeJ
	for k:= 0; k<len(services.Services); k++{
		go func(k int, services Services){
			exeJ(services, k, alertFileEdited, path, services.Services[k].Name)
		}(k, services)
	}
	//launch onModify
	go func(){
		watchFile(path, alertFileEdited)
	}()
	//wait for all exeJ to be completed
	wg.Wait()

}

//takes services struct, services.length, onModify, json for settings path, service name
func exeJ(services Services, number int, alertFileEdited chan string, path string, name string){
	var mode string = services.Services[number].Mode
	var nfqCoonfig uint16 = uint16(services.Services[number].Nfq)
	var regex  = strings.Join(services.Services[number].RegexList,"|")
	var protocol = services.Services[number].Protocol

	fmt.Println("Services -> ", name)
	// fmt.Println("Regex -> ", regex)
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
	reg, _ := regexp.Compile(regex)
	//function executed for every packet in input
	fn := func(a nfqueue.Attribute) int {
		select {
			// if the json is updated, update the regex
			case <- alertFileEdited:
				services = readJson(path) 
				regex = strings.Join(services.Services[number].RegexList,"|")
				reg, _ = regexp.Compile(regex)
			default:
		}
		
		//take packet id
		id := *a.PacketID
		//allocate byte array for packet payload
		payload := make([]byte, len(*a.Payload))
		//copy packet payload to payload variable
		copy(payload, *a.Payload)
		//payload var stringify()
		payloadString := string(payload)
		
		//calculate offset for ignore IP and TCP/UDP headers
		var offset int
		if (protocol == "udp"){
			offset = 20 + 8
		} else if (protocol == "tcp"){
			offset = 20 + ((int(payload[32:33][0]) >> 4)*32)/8
		}

		if(mode == "b"){ //blacklist (if there is a match with the regex, drop the packet)
			if (len(regex) == 0){
				nf.SetVerdict(id, nfqueue.NfAccept)
			} else if(reg.MatchString(payloadString[offset:])){
				log.Print("DROP", " services -> ", name)
				nf.SetVerdict(id, nfqueue.NfDrop)
			} else {
				// log.Print("ACCEPT", " services -> ", name)
				nf.SetVerdict(id, nfqueue.NfAccept)
			}
			// fmt.Printf("%x\n", payloadString[offset:])
		}else if (mode == "w"){ //whitelist (if there is a match with the regex, accept the packet)
			if(reg.MatchString(payloadString[offset:])){
				// log.Print("ACCEPT", " services -> ", name)
				nf.SetVerdict(id, nfqueue.NfAccept)
			} else {
				log.Print("DROP", " services -> ", name)
				nf.SetVerdict(id, nfqueue.NfDrop)
			}
			// fmt.Printf("%x\n", payloadString[offset:])
		}
		return 0
	}

	r := func(e error) int {
		fmt.Println("Error")
		return 0
	}
	//add to nfqueue callback fn for every packet that matches the rules
	err = nf.RegisterWithErrorFunc(ctx, fn, r)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Block till the context expires
	<-ctx.Done()
}

func exeC(mode string, nfqCoonfig uint16, regex string, protocol string){
	fmt.Println("Regex -> ", regex)
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
	reg, _ := regexp.Compile(regex)

	fn := func(a nfqueue.Attribute) int {
		id := *a.PacketID
		payload := make([]byte, len(*a.Payload))
		copy(payload, *a.Payload)
		payloadString := string(payload)
		
		//calculate offset for ignore IP and TCP/UPD headers
		var offset int
		if (protocol == "udp"){
			offset = 20 + 8
		} else if (protocol == "tcp"){
			offset = 20 + ((int(payload[32:33][0]) >> 4)*32)/8
		}

		if(mode == "b"){ //blacklist (if there is a match with the regex, drop the packet)
			if(reg.MatchString(payloadString[offset:])){
				//TODO: do not drop, modify flag
				nf.SetVerdict(id, nfqueue.NfDrop)
			} else {
				nf.SetVerdict(id, nfqueue.NfAccept)
			}
			fmt.Printf("%x\n", payloadString[offset:])
		}else if (mode == "w"){ //whitelist (if there is a match with the regex, accept the packet)
			if(reg.MatchString(payloadString[offset:])){
				nf.SetVerdict(id, nfqueue.NfAccept)
			} else {
				//TODO: do not drop, modify flag
				nf.SetVerdict(id, nfqueue.NfDrop)
			}
			fmt.Printf("%x\n", payloadString[offset:])
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
	var regexFlag = flag.String("r", "", "Regex to match, follow this format: '(regex1)|(regex2)|...'")
	flag.Parse()

	//some change for the flags
	inType := *inTypeFlag
	nfqCoonfig := uint16(*nfqFlag)
	mode := *modeFlag
	protocol := *protocolFlag
	port := *dportFlag
	path := *pathFlag
	regex := *regexFlag

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
		exeC(mode, nfqCoonfig, regex, protocol)
	}

}
