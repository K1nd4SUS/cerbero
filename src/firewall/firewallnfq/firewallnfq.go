package firewallnfq

import (
	"cerbero3/firewall/headers"
	"cerbero3/firewall/rules"
	"cerbero3/logs"
	"cerbero3/metrics"
	"cerbero3/services"
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/florianl/go-nfqueue"
)

// errors that are not relevant to the functionality
// of the firewall
var skippableErrors = map[string]bool{
	"netlink receive: recvmsg: no buffer space available": true,
}

func StartFirewallForService(rr rules.RemoveRules, serviceIndex int) {
	nfqConfig := nfqueue.Config{
		NfQueue:      services.Services[serviceIndex].Nfq,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		// TODO: check this value
		WriteTimeout: 15 * time.Millisecond,
	}

	logs.PrintDebug(fmt.Sprintf(`Opening nfq for service "%v"...`, services.Services[serviceIndex].Name))
	nfq, err := nfqueue.Open(&nfqConfig)
	if err != nil {
		logs.PrintCritical(err.Error())
		rr <- true
		return
	}
	defer nfq.Close()
	logs.PrintDebug(fmt.Sprintf(`Opened nfq for service "%v".`, services.Services[serviceIndex].Name))

	logs.PrintDebug(fmt.Sprintf(`Starting background task for service "%v"...`, services.Services[serviceIndex].Name))
	ctx := context.Background()
	err = nfq.RegisterWithErrorFunc(ctx, func(packet nfqueue.Attribute) int {
		defer func() {
			if err := recover(); err != nil {
				logs.PrintError(fmt.Sprintf("%v", err))

				// if handling the packet panics the program, then
				// we drop it immediately
				nfq.SetVerdict(*packet.PacketID, nfqueue.NfDrop)
			}
		}()

		return handlePacket(nfq, &packet, serviceIndex)
	}, func(err error) int {
		// checks if the error is in the list
		if skippableErrors[err.Error()] {
			logs.PrintError(err.Error())
			return 0
		}

		logs.PrintCritical(err.Error())
		rr <- true
		// we are inside a lambda function, if we do "return" the nfq
		// will not stop. We need to close the nfq directly; it will
		// then do "nfq.Close()" from the defer above, but it will
		// just return an error and ignore it
		nfq.Close()

		// this is a signal to stop receiving messages:
		// https://pkg.go.dev/github.com/florianl/go-nfqueue#ErrorFunc
		return 1
	})
	if err != nil {
		logs.PrintCritical(err.Error())
		rr <- true
		return
	}
	logs.PrintDebug(fmt.Sprintf(`Started background task for service "%v".`, services.Services[serviceIndex].Name))

	// block this thread until nfq is done handling the packets,
	// basically until forever
	<-ctx.Done()
}

func handlePacket(nfq *nfqueue.Nfqueue, packet *nfqueue.Attribute, serviceIndex int) int {
	payload := make([]byte, len(*packet.Payload))
	copy(payload, *packet.Payload)

	var offset int
	if services.Services[serviceIndex].Protocol == "udp" {
		offset = headers.GetUDPHeaderLength()
	} else if services.Services[serviceIndex].Protocol == "tcp" {
		offset = headers.GetTCPHeaderLength(payload)
	}
	payloadString := string(payload)[offset:]

	var droppedRegex string
	verdict := nfqueue.NfAccept
	for _, regex := range services.Services[serviceIndex].RegexList {
		matcher := regexp.MustCompile(regex)
		if matcher.MatchString(payloadString) {
			verdict = nfqueue.NfDrop
			droppedRegex = regex
			break
		}
	}
	nfq.SetVerdict(*packet.PacketID, verdict)
	isDropped := verdict == nfqueue.NfDrop

	metrics.IncrementService(serviceIndex, isDropped)
	if isDropped {
		metrics.IncrementRegex(droppedRegex)
	}
	logs.PrintDebug(fmt.Sprintf(`"%v": %v packet %q.`, services.Services[serviceIndex].Name, func() string {
		if verdict == nfqueue.NfAccept {
			return "accepted"
		} else if isDropped {
			return "dropped"
		}
		return ""
	}(), payloadString))

	// this is a signal to keep receiving messages:
	// https://pkg.go.dev/github.com/florianl/go-nfqueue#ErrorFunc
	return 0
}
