package services

import (
	"bufio"
	"bytes"
	"cerbero/configuration"
	"cerbero/logs"
	"cerbero/random"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dlclark/regexp2"
	"io"
	"math"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/coreos/go-iptables/iptables"
	"github.com/fsnotify/fsnotify"
)

const (
	serviceNameRegex = `^[a-z0-9_]*$`
)

var socketInitializationString = ""

func Configure(config configuration.Configuration) {
	if !configuration.IsCerberoSocketSet(config) {
		return
	}

	socketInitializationString = random.String(
		random.LowerAlphabet+random.UpperAlphabet+random.DigitAlphabet+random.SpecialAlphabet,
		16,
	)
	logs.PrintInfo(fmt.Sprintf(`Using the following Cerbero socket initialization string: "%v".`, socketInitializationString))
}

var Services []Service

type Service struct {
	Nfq       uint16   `json:"nfq"`
	Name      string   `json:"name"`
	Protocol  string   `json:"protocol"`
	Port      int      `json:"port"`
	Chain     string   `json:"chain"`
	RegexList []string `json:"regexes"`
	Matchers  []*regexp2.Regexp
}

func LoadConfigFile(configurationFile string) error {
	logs.PrintDebug(fmt.Sprintf(`Reading configuration file from "%v"...`, configurationFile))
	b, err := os.ReadFile(configurationFile)
	if err != nil {
		return err
	}
	logs.PrintDebug(fmt.Sprintf(`Read configuration file from "%v".`, configurationFile))

	if err = LoadJSON(b); err != nil {
		return err
	}

	logs.PrintDebug("Starting thread listening for file updates in the background...")
	if err = watchForConfigFileChanges(configurationFile); err != nil {
		return err
	}

	return nil
}

func watchForConfigFileChanges(configurationFile string) error {
	logs.PrintDebug("Creating file watcher...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	logs.PrintDebug("Created file watcher.")
	// do NOT close the watcher at the end of the function,
	// it has to keep working forever
	// defer watcher.Close()

	logs.PrintDebug(fmt.Sprintf(`Adding file "%v" to watcher...`, configurationFile))
	err = watcher.Add(configurationFile)
	if err != nil {
		return err
	}
	logs.PrintDebug("Added file to watcher")

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Has(fsnotify.Write) {
					logs.PrintInfo("Detected configuration file change.")
					b, err := os.ReadFile(configurationFile)
					if err != nil {
						logs.PrintError(fmt.Sprintf("Error while reading the file: %v.", err.Error()))
						continue
					}

					// the file may still be updating so an empty
					// file is read instead of the actual file.
					// the actual one will be read in the next
					// iteration anyways
					if len(b) == 0 {
						continue
					}

					if err = LoadJSON(b); err != nil {
						logs.PrintError(fmt.Sprintf("Error while loading JSON file: %v.", err.Error()))
					}
				}

			case err := <-watcher.Errors:
				if err != nil {
					logs.PrintError(fmt.Sprintf("Error while watching for file changes: %v.", err.Error()))
				}
			}
		}
	}()

	return nil
}

func LoadCerberoSocket(ip string, port int, attempt int) error {
	logs.PrintDebug(fmt.Sprintf(`Connecting to socket at "%v:%v"...`, ip, port))
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", ip, port))
	if err != nil {
		// if it's the first time starting the firewall, the attempt
		// is going to be equal to 0; any other time, it's going to
		// be greather than 0; we only want to attempt reconnection
		// if the firewall was already running, not if it was just
		// started
		if attempt > 0 {
			waitTime := math.Min(float64(attempt)*2, 30)
			logs.PrintError(fmt.Sprintf("Connection failed. Waiting %v seconds before trying again...", waitTime))
			time.Sleep(time.Duration(waitTime) * time.Second)
			logs.PrintDebug("Attempting reconnection...")
			LoadCerberoSocket(ip, port, attempt+1)
			return nil
		}
		return err
	}
	logs.PrintInfo(fmt.Sprintf(`Connected to socket at "%v:%v".`, ip, port))

	logs.PrintDebug(`Sending first byte to socket...`)
	conn.Write([]byte(socketInitializationString))

	logs.PrintDebug("Waiting for the first configuration file from socket...")
	if err = waitForCerberoSocketJSON(conn, true); err != nil {
		return err
	}

	logs.PrintDebug("Starting thread listening for file updates in the background...")
	go func() {
		for {
			if err = waitForCerberoSocketJSON(conn, false); err != nil {
				if err == io.EOF {
					logs.PrintInfo("Socket disconnected. Attempting reconnection...")
					LoadCerberoSocket(ip, port, attempt+1)
					break
				}
				logs.PrintError(fmt.Sprintf("Error while waiting for the next Cerbero socket JSON: %v.", err.Error()))
			}
		}
	}()

	return nil
}

func waitForCerberoSocketJSON(conn net.Conn, firstFile bool) error {
	reader := bufio.NewReader(conn)
	b64, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	if firstFile {
		logs.PrintDebug("Got the first configuration file from socket.")
	} else {
		logs.PrintInfo("Got configuration file from socket.")
	}

	logs.PrintDebug("Decoding the configuration file base64 to normal JSON...")
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	logs.PrintDebug("Decoded the configuration file base64 to normal JSON...")

	err = LoadJSON(b)
	if err != nil {
		return err
	}

	return nil
}

func LoadJSON(b []byte) error {
	logs.PrintDebug("Parsing JSON file...")
	err := json.Unmarshal(b, &Services)
	if err != nil {
		return err
	}
	logs.PrintDebug(func() string {
		buffer := &bytes.Buffer{}
		// this should have been already parsed by the Unmarshal
		// function, if it outputs an error there's clearly something
		// wrong with the compiler or something else
		json.Compact(buffer, b)

		return fmt.Sprintf("Parsed JSON file: %v", buffer.String())
	}())

	CompileMatchers()

	return nil
}

func CheckServicesValues() error {
	// we initialize iptables in order to not initialize
	// it multiple times in the loop
	ipt, err := iptables.New()
	if err != nil {
		return errors.New(fmt.Sprintf("Error while initializing iptables: %v.", err.Error()))
	}

	for _, service := range Services {
		if matches, _ := regexp.MatchString(serviceNameRegex, service.Name); !matches {
			return errors.New(fmt.Sprintf(`Invalid name for service "%v", must match the regex %v.`, service.Name, serviceNameRegex))
		}

		if !(1 <= service.Nfq && service.Nfq <= 65535) {
			return errors.New(fmt.Sprintf(`Invalid nfq for service "%v", must be a value from 1 to 65535.`, service.Name))
		}

		if !(service.Protocol == "tcp" || service.Protocol == "udp") {
			return errors.New(fmt.Sprintf(`Invalid protocol for service "%v", must be set to either "tcp" or "udp".`, service.Name))
		}

		if !(1 <= service.Port && service.Port <= 65535) {
			return errors.New(fmt.Sprintf(`Invalid port for service "%v", must be a value from 1 to 65535.`, service.Name))
		}

		doesChainExist, err := ipt.ChainExists("filter", service.Chain)
		if err != nil {
			return errors.New(fmt.Sprintf("Error while checking if chain exists: %v.", err.Error()))
		}
		if !doesChainExist {
			return errors.New(fmt.Sprintf("The given chain does not exist: %v.", service.Chain))
		}
	}

	return nil
}

func CompileMatchers() {
	logs.PrintDebug("Compiling regex matchers...")

	for index := range Services {
		Services[index].Matchers = nil
		for _, regex := range Services[index].RegexList {
			// TODO: support regex options
			re, err := regexp2.Compile(regex, 0)
			if err != nil {
				logs.PrintError(fmt.Sprintf(`Error parsing regex "%v" for service %v: %v.`, regex, Services[index].Name, err))
				continue
			}

			Services[index].Matchers = append(Services[index].Matchers, re)
		}
	}

	logs.PrintDebug("Compiled regex matchers.")
}
