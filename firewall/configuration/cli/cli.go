package cli

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	File struct {
		FilePath string `name:"path" arg:"" help:"Relative or absolute path to the JSON configuration file."`
	} `name:"file" cmd:"" help:"Use Cerbero with a local file."`
	Socket struct {
		Address string `name:"address" arg:"" help:"The server to which Cerbero will connect to update the configuration file. The format must be <ip>:<port>."`
	} `name:"socket" cmd:"" help:"Use Cerbero with a remote custom server or with the Cerbero web plugin."`

	MetricsPort int    `name:"metrics-port" short:"P" help:"Port used for the metrics server." default:"9090"`
	ColoredLogs bool   `name:"colored-logs" short:"c" help:"Enable colors for logs. They will not appear in the logfile." default:"false"`
	LogFile     string `name:"log-file" short:"L" help:"File used to output logs." default:"/var/log/cerbero/status.log"`
	Verbose     bool   `name:"verbose" short:"v" help:"Enable DEBUG-level logging." default:"false"`
}

func Parse() CLI {
	flags := CLI{}
	kong.Parse(&flags)

	return flags
}
