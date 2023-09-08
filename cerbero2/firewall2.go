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
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"regexp"
	"strconv"
	"sync"
	"time"

	//"fmt"

	//gopacket "github.com/google/gopacket"
	//layers "github.com/google/gopacket/layers"
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

// stats structs
type Stats struct {
	FileEdits     uint32
	ServiceAccess []ServiceAccess // there will be a service access for each service
}

type ServiceAccess struct {
	Service Service // containing useful info such as the service name and port
	Hits    []Hit   // list of hits registered on that particular service
}

type Hit struct {
	Resource string // a hit is characterized by a hit resource, the method used and number of accesses (and blocked ones)
	Method   string
	Counter  uint64
	Blocked  uint64
}

// when starting the firewall a new serviceAccess item is added for each registered service (from config.json). Then, when receiving a request, a check is made to verify that it is a new access resource. In this case, a new item in Hit is added. If not,
// the "Counter" is just increased. Finally, based on the verdict, "Blocked" may be increased also.

var stats Stats

//map
type ResInfo struct {
	Accessed     bool
	Idx 		 int 
}

//! RAGIONAMENTO
//* 1. la mappa ha un campo aggiuntivo: il timestamp/unix time dell'ultima volta che è stato usato quel pacchetto 
//* 2. ogni volta che ci risulta quel pacchetto, aggiorniamo il tempo
//* 3. una goroutine runna una volta ogni tot secondi, controlla per ogni pacchetto se non si sentono sue notizie da più del delta fissato, e a quel punto lo elimina

//flag for chain selection on iptables
var chainType="DOCKER-INGRESS"

// logs
var warnings = make(chan string, 1)
var normal = make(chan string, 1)
var infos = make(chan string, 1)
var success = make(chan string, 1)

var newStats = make(chan uint8) // to know when a new stats (to be printed on the cli) is ready

func printWarnings() {
	for msg := range warnings {
		log.Println("\x1b[38;5;202m\t" + msg + "\033[0m") // orange
	}
}
func printSuccess() {
	for msg := range success {
		log.Println("\x1b[38;5;10m\t" + msg + "\033[0m") // green
	}
}
func printInfos() {
	for msg := range infos {
		log.Println("\x1b[38;5;51m\t" + msg + "\033[0m") // cyan
	}
}
func printErrors(msg string) {
	log.Println("\x1b[38;5;1m\t" + msg + "\033[0m") // red
	os.Exit(127)
}
func printNormal() {
	for msg := range normal {
		log.Println("\t" + msg)
	}
}

func printStats() { // just printing stats on the cli
	for {
		<-newStats
		log.Printf("\x1b[47;5;1m\t%+v\033[0m", stats) // that one
	}
}

// serialize input
func readJson(path string) Services {
	jsonFile, _ := os.Open(path)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var services Services
	if json.Valid(byteValue) {
		json.Unmarshal(byteValue, &services)
		stats.FileEdits++
		newStats <- 1
		return services
	}
	warnings <- "An error was found in the config file!"
	var noServices Services
	return noServices
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
func watchFile(path string, alertFile chan string) {
	oldHash := hash(path)
	for {
		time.Sleep(5 * time.Second)
		newHash := hash(path)
		if oldHash != newHash {
			infos <- "Configuration file edited"
			alertFile <- "-"
		}
		oldHash = newHash
	}
}

// check params validity
func checkParams(serv *Service, nfqConfig uint16) {

	// for every param, if param is not allowed the execution is terminated, else everything can go on

	//checks if the procols is correct (must be "tcp" or "udp")
	if serv.Protocol != "tcp" && serv.Protocol != "udp" {
		printErrors("Invalid argument for flag -p, must be set to 'tcp' or 'udp'")
	}

	//check if the port number is right
	if serv.Dport < 1 || serv.Dport > 65535 {
		printErrors("Invalid argument for flag -dport, the value need to be between 1 and 65535")
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
		printErrors("Invalid argument for flag -nfq, the value need to be between 1 and 65535")
	}

	// control if file exists
	_, err := os.Open(path)
	if err != nil { //if it doesn't
		printErrors("File not found") //close.
	}
	//everything is fine, the file is there

	services := readJson(path)
	if len(services.Services) == 0 {
		printErrors("No services or error in config!")
	}

	infos <- "services parsed"

	for k := 0; k < len(services.Services); k++ {
		checkParams(&services.Services[k], (nfqConfig + uint16(k)))
	}

	return services
}



// apply filters and keep rules updated
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

	var lastResourceIndex int // to manage fragments. It stores the index in stats of the resource the packets is asking for. This is useful because if fwfilter is processing an intermediate fragment, it knows which resource must be increased in accesses (and maybe blocks) in stats
	var fragMap = make(map[string]ResInfo)
	
	//se non riesce ad aprire il socket chiudi
	nf, err := nfqueue.Open(&config)
	if err != nil {
		printErrors("could not open nfqueue socket")
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
	fn := func(packet nfqueue.Attribute) int {
		select {
		// if the json is updated, update the regex
		case <-alertFileEdited:

			tempServices := readJson(path)
			//if json is valid,then apply changes
			if len(tempServices.Services) != 0 {
				services = tempServices
			}
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
		id := *packet.PacketID

		// // Decode a packet
		//packet := gopacket.NewPacket(*a.Payload, layers.LayerTypeEthernet, gopacket.Default)
		// // Get the TCP layer from this packet
		// var tcpLayer = packet.Layer(layers.LayerTypeTCP);
		// if  tcpLayer != nil {
		// log.Println("This is a TCP packet!")
		// // Get actual TCP data from this layer
		// tcp, _ := tcpLayer.(*layers.TCP)
		// log.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
		// }
		// log.Println(tcpLayer)

		notManaged := true

		//allocate byte array for packet payload
		payload := make([]byte, len(*packet.Payload))

		//log.Println(len(*packet.Payload))
		//copy packet payload to payload variable
		copy(payload, *packet.Payload)

		//payload var stringify()
		payloadString := string(payload)

		//calculate offset for ignore IP and TCP/UDP headers
		var offset int
		if protocol == "udp" {
			offset = 20 + 8
		} else if protocol == "tcp" {
			offset = 20 + ((int(payload[32:33][0])>>4)*32)/8
		}
		fmt.Println("\x1b[38;5;129m","INIZIO PACCHETTO","\033[0m")
		log.Println("lunghezza ", len(payloadString[offset:]), " offset ", offset)
		var newResource string
		var newMethodType = ""
		var splitted string
		methReg, _ := regexp.Compile("(GET )|(POST )|(PUT )|(PATCH )|(DELETE )|(HEAD )|(CONNECT )|(OPTIONS )|(TRACE )")
		boundReg, _ := regexp.Compile("boundary=------------------------")
		bound2Reg, _ := regexp.Compile("--------------------------")
		var boundary = []string{""}
		
		if len(payloadString[offset:]) > 0 { // if the packet contains anything, TODO: exploitable to cause crashes
			newResource = methReg.Split(payloadString[offset:], 1)[0]
			//log.Println("\n\n\n\n", newResource)
			fmt.Println(newResource)
			splitted = strings.Split(newResource, "HTTP")[0]
			
			newResource = strings.Split(splitted," ")[1]
			newMethodType = strings.Split(splitted," ")[0]
			//newResource = strings.Split(payloadString[offset+4:],  reg)[0] // retrieving the resource name
			
		 } else {
			newResource = methReg.Split(payloadString[offset:], 1)[0]
			newResource = strings.Split(newResource, "HTTP")[0]
		}

		if methReg.MatchString(payloadString[offset:]) { // if this packet is a fragment and is the first fragment, it must contain GET/POST string (can be improved  maybe including also the HTTP/*.* string)
			
			//* take the boundary, if the regex is found, then the array MUST be len > 1
			boundary = boundReg.Split(payloadString[offset:], -1)			
			
			alreadyAccessed := false
			var i int
			if len(stats.ServiceAccess[number].Hits) > 0 {
				for i = 0; i < len(stats.ServiceAccess[number].Hits) && !alreadyAccessed; i++ { // looking for the already accessed resource
					if (stats.ServiceAccess[number].Hits[i].Resource == newResource) && (stats.ServiceAccess[number].Hits[i].Method == newMethodType){
						alreadyAccessed = true
						fmt.Println("INDICE RISORSA",stats.ServiceAccess[number].Hits[i].Resource," IDX",i)
					}
				}
			}

			if i > 0 && (i != len(stats.ServiceAccess[number].Hits) || alreadyAccessed) {
				i--
			}
			lastResourceIndex = i
			
			//!TODO con alcuni metodi GET il lastresidx non cambia!!
			if !alreadyAccessed { // creating a new accessed resource in stats
				var newHit Hit
				newHit.Resource = newResource
				newHit.Method = newMethodType
				newHit.Counter++
				stats.ServiceAccess[number].Hits = append(stats.ServiceAccess[number].Hits, newHit)
			} else {
				stats.ServiceAccess[number].Hits[lastResourceIndex].Counter++
			}
			fmt.Println("IDX ",i," len", len(stats.ServiceAccess[number].Hits))

			
			//* insert in the map key= boundary & value=0 IF the key was found
			if(len(boundary) > 1){
				var resStruct ResInfo
				resStruct.Accessed = false
				resStruct.Idx = lastResourceIndex
				
				fragMap[boundary[1][:16]]=resStruct
			}
			
		}

		// if this conditional is not entered, it means that we are handling an intermediate fragment, so only a verdict must be given (surely we are not adding a new resource access)
		hexReg, _ := regexp.Compile("^[0-9a-fA-F]+$")
		
		if hasWhitelist { //whitelist (if there is a match with the regex, accept the packet)

			if !whitelistMatcher.Contains([]byte(payloadString[offset:])) {
				warnings <- "packet dropped because whitelist " + services.Services[number].Name // + "ID: " + strconv.FormatUint(uint64(id), 10)
				nf.SetVerdict(id, nfqueue.NfDrop)

				boundary = bound2Reg.Split(payloadString[offset:], -1) 

				//* we use shortcircuiting for avoiding a crash here, we check if the fragment is not already counted then we update is value
				if len(boundary) > 1 && !fragMap[boundary[1][:16]].Accessed {
					stats.ServiceAccess[number].Hits[fragMap[boundary[1][:16]].Idx].Blocked++
					tempStruct := fragMap[boundary[1][:16]]
					tempStruct.Accessed = true
					fragMap[boundary[1][:16]] = tempStruct
					
				}else if(len(boundary) == 1 && !hexReg.MatchString(boundary[0][:16])){
					
					stats.ServiceAccess[number].Hits[lastResourceIndex].Blocked++
				
				}
				notManaged = false
			}
		}

		if hasBlacklist && notManaged { //blacklist (if there is a match with the regex, drop the packet)

			if blacklistMatcher.Contains([]byte(payloadString[offset:])) {
				warnings <- "packet dropped because of " + services.Services[number].Name + " blacklist" // + "ID: " + strconv.FormatUint(uint64(id), 10)
				nf.SetVerdict(id, nfqueue.NfDrop)
				boundary = bound2Reg.Split(payloadString[offset:], -1)
				
				//* we use shortcircuiting for avoiding a crash here, we check if the fragment is not already counted then we update is value
				if len(boundary) > 1 && !fragMap[boundary[1][:16]].Accessed {

					stats.ServiceAccess[number].Hits[fragMap[boundary[1][:16]].Idx].Blocked++
					tempStruct := fragMap[boundary[1][:16]]
					tempStruct.Accessed = true
					fragMap[boundary[1][:16]] = tempStruct
				}else if(len(boundary) == 1 && !hexReg.MatchString(boundary[0][:16]) ){	
					
					stats.ServiceAccess[number].Hits[lastResourceIndex].Blocked++
				
				}
				notManaged = false
			}
		}

		if notManaged {
			nf.SetVerdict(id, nfqueue.NfAccept)
		}
		newStats <- 1
		warnings <- payloadString[offset:] // just printing the payload - CONCURRENT <- is printed after "FINE PACCHETTO"

		return 0
	}

	r := func(e error) int {
		printErrors("Error")
		return 42
	}

	//add to nfqueue callback fn for every packet that matches the rules
	err = nf.RegisterWithErrorFunc(ctx, fn, r)
	if err != nil {
		log.Println(err)
		return
	}

	// Block until the context expires
	<-ctx.Done()
}

// set rules on iptables and call fwFilter
func setRules(services Services, path string) {
	for _, ser := range services.Services {
		log.Println(ser)
	}

	//loop for create iptables rules
	for k := 0; k < len(services.Services); k++ {
		cmd := exec.Command("iptables", "-I", chainType, "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
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
			cmd := exec.Command("iptables", "-D", chainType, "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
			cmd.Run()
		}
		log.Println("\x1b[38;5;10m\tDone!\033[0m")
		os.Exit(0)
	}()

	//start waitgroup
	var wg sync.WaitGroup

	//onmodify for json
	alertFileEdited := make(chan string, 10)

	//create waitgroup
	wg.Add(len(services.Services) + 1)

	//loop for start the go routines with fwFilter
	for k := 0; k < len(services.Services); k++ {
		go func(k int, services Services) {
			fwFilter(services, k, alertFileEdited, path)
		}(k, services)

		var newServiceAccess ServiceAccess // adding a new service access in stats for each service from the config file
		newServiceAccess.Service = services.Services[k]
		stats.ServiceAccess = append(stats.ServiceAccess, newServiceAccess)
		newStats <- 1
	}

	//launch onModify
	go watchFile(path, alertFileEdited)

	//wait for all fwFilter to be completed
	wg.Wait()

}

func statsHandler(w http.ResponseWriter, r *http.Request) { // handling stats queries over API
	if r.URL.Path != "/metrics" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	marshaled, err := json.MarshalIndent(stats, "", "   ")
	if err != nil {
		log.Fatalf("marshaling error: %s", err)
	}
	fmt.Fprintf(w, string(marshaled)) // sending stats in pretty JSON

}

func main() {

	go printWarnings()
	go printNormal()
	go printInfos()
	go printSuccess()

	go printStats()

	http.HandleFunc("/metrics", statsHandler) // giving stats on /stats :8082
	go http.ListenAndServe(":8082", nil)

	success <- "Service started"

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
	//chain specification
	var dockerized = flag.String("docker", "y", "Are the services on docker? [Y/n]")
	
	flag.Parse()
	success <- "Flags parsed"

	if *dockerized == "n" {
		chainType = "INPUT"
	}
	
	infos <- "chain "+chainType+" selected"
	nfqConfig := uint16(*nfqFlag)
	path := *pathFlag

	//checks flags
	serviceList := checkIn(path, nfqConfig)

	//here we will call a func that executes everything
	setRules(serviceList, path)

}
