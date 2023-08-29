package main

import (
	//"context"
	//"crypto/sha256"
	//"encoding/hex"
	"encoding/json"
	"flag"
	//"fmt"
	//"io"
	"io/ioutil"
	"log"
	"os"
	//"os/exec"
	//"os/signal"
	//"regexp"
	//"strconv"
	//"strings"
	//"sync"
	//"time"

	//nfqueue "github.com/florianl/go-nfqueue"
)

// structs

type Services struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name		string 		`json:"name"`
	Nfq 		uint16 		
	Mode 		string 		`json:"mode"`
	Protocol	string 		`json:"protocol"`
	Dport 		int 		`json:"dport"`
	RegexList	[]string	`json:"regexList"`
}

//serialyze input
func readJson(path string)(Services){
	jsonFile, _ := os.Open(path)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var services Services
	json.Unmarshal(byteValue, &services)
	return services
}

//check params validity
func checkParams(serv *Service, nfqConfig uint16){

	// for every param, if param is not allowed the execution is terminated, else everything can go on
	
	//check if mode is allowed (must be "w" or "b")
	if(serv.Mode != "w" && serv.Mode != "b"){
		log.Println("Invalid argument for flag -mode, must be set to 'w' or 'b'")
		os.Exit(127)
	}

	//checks if the procols is correct (must be "tcp" or "udp")
	if(serv.Protocol != "tcp" && serv.Protocol != "udp"){
		log.Println("Invalid argument for flag -p, must be set to 'tcp' or 'udp'")
		os.Exit(127)
	}

	//check if the port number is right
	if(serv.Dport < 1 || serv.Dport > 65535){
		log.Println("Invalid argument for flag -dport, the value need to be between 1 and 65535")
		os.Exit(127)
	}
	
	//assigning nfq id
	serv.Nfq = nfqConfig

}

//load params
func checkIn(path string, nfqConfig uint16)(Services){
	
	/*
		EDITS:
			- removed nfq number -> we'll insert them manually
			- removed cli config -> only json allowed in 21st century
	*/

	// check nfq number
	if(nfqConfig < 1 || nfqConfig > 65535){
		log.Println("Invalid argument for flag -nfq, the value need to be between 1 and 65535")
		os.Exit(127)
	}

	// control if file exists
	_, err := os.Open(path)
	if (err != nil){	//if it doesn't
		log.Println("File not found") //print
		os.Exit(127)	//close.
	}
	//everything is fine, the file is there
	
	services := readJson(path)
	
	for k:= 0; k<len(services.Services); k++{
		checkParams(&services.Services[k],(nfqConfig+uint16(k)))
	}
	
	return services
}

func execJson(services Services){
	for _,ser := range services.Services{
		log.Println(ser)
	}
	// services := readJson(path)
	// //loop for create iptables rules
	// for k:= 0; k<len(services.Services); k++{
	// 	cmd := exec.Command("iptables", "-I", "INPUT", "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
	// 	cmd.Run()
	// }
	// //prepare oninterrupt event
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// go func(){
	// 	<-c
	// 	log.Println("\nRemoving iptables rule")
	// 	//loop for delete iptables rules
	// 	for k:= 0; k<len(services.Services); k++{
	// 		cmd := exec.Command("iptables", "-D", "INPUT", "-p", services.Services[k].Protocol, "--dport", strconv.FormatInt(int64(services.Services[k].Dport), 10), "-j", "NFQUEUE", "--queue-num", strconv.FormatInt(int64(services.Services[k].Nfq), 10))
	// 		cmd.Run()
	// 	}
	// 	log.Println("Done!")
	// 	os.Exit(0)
	// }()	
	// //start waitgroup 
	// var wg sync.WaitGroup
	// //onmodify for json
	// alertFileEdited := make(chan string)
	// wg.Add(len(services.Services)+1)
	// //loop for start the go routines with exeJ
	// for k:= 0; k<len(services.Services); k++{
	// 	go func(k int, services Services){
	// 		exeJ(services, k, alertFileEdited, path, services.Services[k].Name)
	// 	}(k, services)
	// }
	// //launch onModify
	// go func(){
	// 	watchFile(path, alertFileEdited)
	// }()
	// //wait for all exeJ to be completed
	// wg.Wait()

}


func main() {

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
	
	nfqConfig := uint16(*nfqFlag)
	path := *pathFlag

	//checks flags
	serviceList := checkIn(path, nfqConfig)

	//here we will call a func that executes everything
	execJson(serviceList)
	

}

