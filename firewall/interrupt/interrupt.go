package interrupt

import (
	"cerbero/firewall/rules"
	"os"
	"os/signal"
)

func Listen(rr rules.RemoveRules) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		// wait for a stop signal to arrive
		<-stop
		rr <- true
	}()
}
