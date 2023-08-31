package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"

	//"regexp"
	"strconv"
	//"strings"
	"sync"
	"time"

	//"fmt"

	ahocorasick "github.com/cloudflare/ahocorasick"
	nfqueue "github.com/florianl/go-nfqueue"
)

// structs

type Services struct {
	Services []Service `json:"services"`
}

type Rule struct {
	Type    string   `json:"type"`
	Filters []string `json:"filters"`
}

type Rules struct {
	Blacklist []Rule `json:"blacklist"`
	Whitelist []Rule `json:"whitelist"`
}

type Service struct {
	Name      string `json:"name"`
	Nfq       uint16
	Protocol  string `json:"protocol"`
	Dport     int    `json:"dport"`
	RulesList Rules  `json:"rulesList"`
}

// logs
var warnings = make(chan string, 1)
var errors = make(chan string, 1)
var normal = make(chan string, 1)
var infos = make(chan string, 1)
var success = make(chan string, 1)

func printWarnings() {
	for msg := range warnings {
		log.Println("\x1b[38;5;202m\t" + msg + "\033[0m") // orange
	}
}
func printSuccess() {
	for msg := range success {
		log.Println("\x1b[38;5;10m\t" + msg + "\033[0m") // orange
	}
}
func printInfos() {
	for msg := range infos {
		log.Println("\x1b[38;5;51m\t" + msg + "\033[0m") // cyan
	}
}
func printErrors() {
	for msg := range errors {
		log.Println("\x1b[38;5;1m\t" + msg + "\033[0m") // red
	}
}
func printNormal() {
	for msg := range normal {
		log.Println("\t" + msg)
	}
}

// serialize input
func readJson(path string) Services {
	jsonFile, _ := os.Open(path)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var services Services
	json.Unmarshal(byteValue, &services)
	return services
}

// hash a string, given file path
func hash(path string) (hash string) {
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

// onModify
func watchFile(path string, canale chan string) {
	oldHash := hash(path)
	for true {
		time.Sleep(5 * time.Second)
		newHash := hash(path)
		if oldHash != newHash {
			infos <- "Configuration file edited"
			canale <- "-"
		}
		oldHash = newHash
	}
}

// check params validity
func checkParams(serv *Service, nfqConfig uint16) {

	// for every param, if param is not allowed the execution is terminated, else everything can go on

	//checks if the procols is correct (must be "tcp" or "udp")
	if serv.Protocol != "tcp" && serv.Protocol != "udp" {
		log.Println("Invalid argument for flag -p, must be set to 'tcp' or 'udp'")
		os.Exit(127)
	}

	//check if the port number is right
	if serv.Dport < 1 || serv.Dport > 65535 {
		log.Println("Invalid argument for flag -dport, the value need to be between 1 and 65535")
		os.Exit(127)
	}

	//assigning nfq id
	serv.Nfq = nfqConfig

}

// load params
func checkIn(path string, nfqConfig uint16) Services {

	/*
		EDITS:
			- removed nfq number -> we'll insert them manually
			- removed cli config -> only json allowed in 21st century
	*/

	// check nfq number
	if nfqConfig < 1 || nfqConfig > 65535 {
		log.Println("Invalid argument for flag -nfq, the value need to be between 1 and 65535")
		os.Exit(127)
	}

	// control if file exists
	_, err := os.Open(path)
	if err != nil { //if it doesn't
		log.Println("File not found")
		os.Exit(127) //close.
	}
	//everything is fine, the file is there

	services := readJson(path)

	infos <- "services parsed"

	for k := 0; k < len(services.Services); k++ {
		checkParams(&services.Services[k], (nfqConfig + uint16(k)))
	}

	return services
}

// apply filters
func fwFilter(services Services, number int, alertFileEdited chan string, path string) {

	blacklist := services.Services[number].RulesList.Blacklist
	whitelist := services.Services[number].RulesList.Whitelist
	hasBlacklist := (len(blacklist) != 0)
	hasWhitelist := (len(whitelist) != 0)
	var nfqConfig uint16 = uint16(services.Services[number].Nfq)
	var protocol string = services.Services[number].Protocol

	infos <- "activated on: " + services.Services[number].Name

	// Set configuration options for nfqueue
	config := nfqueue.Config{
		NfQueue:      nfqConfig,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 15 * time.Millisecond,
	}

	//se non riesce ad aprire il socket chiudi
	nf, err := nfqueue.Open(&config)
	if err != nil {
		log.Println("could not open nfqueue socket:", err)
		return
	}
	defer nf.Close()

	ctx := context.Background()

	//TODO: controllare sta porcata
	blacklistMatcher := ahocorasick.NewStringMatcher(make([]string, 0))
	if hasBlacklist {
		blacklistMatcher = ahocorasick.NewStringMatcher(blacklist[0].Filters)
	}

	whitelistMatcher := ahocorasick.NewStringMatcher(make([]string, 0))
	if hasWhitelist {
		whitelistMatcher = ahocorasick.NewStringMatcher(whitelist[0].Filters)
	}

	//function executed for every packet in input
	fn := func(a nfqueue.Attribute) int {
		select {
		// if the json is updated, update the regex
		case <-alertFileEdited:

			services = readJson(path)
			blacklist = services.Services[number].RulesList.Blacklist
			whitelist = services.Services[number].RulesList.Whitelist
			hasBlacklist = (len(blacklist) != 0)
			hasWhitelist = (len(whitelist) != 0)
			//blacklistMatcher := ahocorasick.NewStringMatcher(make([]string, 0))
			if hasBlacklist {
				blacklistMatcher = ahocorasick.NewStringMatcher(blacklist[0].Filters)
			}

			//whitelistMatcher := ahocorasick.NewStringMatcher(make([]string, 0))
			if hasWhitelist {
				whitelistMatcher = ahocorasick.NewStringMatcher(whitelist[0].Filters)
			}

		default:
		}

		//take packet id
		id := *a.PacketID

		not_managed := true

		//allocate byte array for packet payload
		payload := make([]byte, len(*a.Payload))

		//copy packet payload to payload variable
		copy(payload, *a.Payload)

		//payload var stringify()
		payloadString := string(payload)

		//calculate offset for ignore IP and TCP/UDP headers
		var offset int
		if protocol == "udp" {
			offset = 20 + 8
		} else if protocol == "tcp" {
			offset = 20 + ((int(payload[32:33][0])>>4)*32)/8
		}

		if hasWhitelist { //whitelist (if there is a match with the regex, accept the packet)

			if !whitelistMatcher.Contains([]byte(payloadString[offset:])) {
				warnings <- "packet dropped " + services.Services[number].Name
				nf.SetVerdict(id, nfqueue.NfDrop)
				not_managed = false
			}
		}

		if hasBlacklist && not_managed { //blacklist (if there is a match with the regex, drop the packet)

			if blacklistMatcher.Contains([]byte(payloadString[offset:])) {
				warnings <- "packet dropped " + services.Services[number].Name
				nf.SetVerdict(id, nfqueue.NfDrop)
				not_managed = false
			}
		}

		if not_managed {
			nf.SetVerdict(id, nfqueue.NfAccept)
		}

		return 0
	}

	r := func(e error) int {
		log.Println("Error")
		return 0
	}

	//add to nfqueue callback fn for every packet that matches the rules
	err = nf.RegisterWithErrorFunc(ctx, fn, r)
	if err != nil {
		log.Println(err)
		return
	}

	// Block till the context expires
	<-ctx.Done()
}

// set rules on iptables and call fwFilter
func setRules(services Services, path string) {
	for _, ser := range services.Services {
		log.Println(ser)
	}

	//loop for create iptables rules
	for k := 0; k < len(services.Services); k++ {
		cmd := exec.Command("iptables", "-I", "INPUT", "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
		cmd.Run()
	}

	//prepare oninterrupt event
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		success <- "Removing iptables rule"
		//loop for delete iptables rules
		for k := 0; k < len(services.Services); k++ {
			cmd := exec.Command("iptables", "-D", "INPUT", "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
			cmd.Run()
		}
		log.Println("Done!")
		os.Exit(0)
	}()

	//start waitgroup
	var wg sync.WaitGroup

	//onmodify for json
	alertFileEdited := make(chan string)

	//create waitgroup
	wg.Add(len(services.Services) + 1)

	//loop for start the go routines with fwFilter
	for k := 0; k < len(services.Services); k++ {
		go func(k int, services Services) {
			fwFilter(services, k, alertFileEdited, path)
		}(k, services)
	}

	//launch onModify
	go func() {
		watchFile(path, alertFileEdited)
	}()

	//wait for all fwFilter to be completed
	wg.Wait()

}

func main() {

	go printErrors()
	go printWarnings()
	go printNormal()
	go printInfos()
	go printSuccess()

	success <- "service started"

	/*
		EDITS:
			- deleted nfqFlag: we insert them manually
			- removed cli config -> only json allowed in 21st century
	*/

	// Send ingoing packets to nfqueue queue 100
	// $ sudo iptables -I INPUT -p tcp --dport 12345 -j NFQUEUE --queue-num 100

	//nfq flag config
	var nfqFlag = flag.Int("nfq", 100, "Queue number (optional, default 100 onwards)")
	//path specification
	var pathFlag = flag.String("path", "./config.json", "Path to the json config file")

	flag.Parse()
	success <- "flags parsed"

	nfqConfig := uint16(*nfqFlag)
	path := *pathFlag

	//checks flags
	serviceList := checkIn(path, nfqConfig)

	//here we will call a func that executes everything
	setRules(serviceList, path)

}
