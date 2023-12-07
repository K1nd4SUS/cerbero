package services

import (
	"bufio"
	"bytes"
	"cerbero3/logs"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	// TODO: change this value
	// the string used by the firewall to communicate
	// to the Cerbero socket that the socket connection
	// has been initialized
	socketInitializedString = ""
)

var Services []Service

type Service struct {
	Nfq       uint16   `json:"nfq"`
	Name      string   `json:"name"`
	Protocol  string   `json:"protocol"`
	Port      int      `json:"port"`
	RegexList []string `json:"regex_list"`
	Matchers  []*regexp.Regexp
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
		if attempt > 0 {
			waitTime := math.Min(float64(attempt)*2, 30)
			logs.PrintError(fmt.Sprintf("Connection failed. Waiting %v seconds before trying again...", waitTime))
			time.Sleep(time.Duration(waitTime) * time.Second)
			logs.PrintInfo("Attempting reconnection...")
			LoadCerberoSocket(ip, port, attempt+1)
			return nil
		}
		return err
	}
	logs.PrintDebug(fmt.Sprintf(`Connected to socket at "%v:%v".`, ip, port))

	logs.PrintDebug(`Sending first byte to socket...`)
	conn.Write([]byte(socketInitializedString))

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
	for _, service := range Services {
		if !(1 <= service.Nfq && service.Nfq <= 65535) {
			return errors.New(fmt.Sprintf(`Invalid nfq for service %v, must be a value from 1 to 65535.`, service.Name))
		}

		if !(service.Protocol == "tcp" || service.Protocol == "udp") {
			return errors.New(fmt.Sprintf(`Invalid protocol for service %v, must be set to either "tcp" or "udp".`, service.Name))
		}

		if !(1 <= service.Port && service.Port <= 65535) {
			return errors.New(fmt.Sprintf(`Invalid port for service %v, must be a value from 1 to 65535.`, service.Name))
		}
	}

	return nil
}

func CompileMatchers() {
	logs.PrintDebug("Compiling regex matchers...")

	for index := range Services {
		Services[index].Matchers = nil
		for _, regex := range Services[index].RegexList {
			Services[index].Matchers = append(Services[index].Matchers, regexp.MustCompile(regex))
		}
	}

	logs.PrintDebug("Compiled regex matchers...")
}
