package firewallnfq

import (
	"cerbero3/firewall/headers"
	"cerbero3/firewall/rules"
	"cerbero3/logs"
	"cerbero3/services"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/florianl/go-nfqueue"
)

func StartFirewallForService(rr rules.RemoveRules, service services.Service) {
	nfqConfig := nfqueue.Config{
		NfQueue:      service.Nfq,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		// TODO: check this value
		WriteTimeout: 15 * time.Millisecond,
	}

	logs.PrintDebug(fmt.Sprintf(`Opening nfq for service "%v"...`, service.Name))
	nfq, err := nfqueue.Open(&nfqConfig)
	if err != nil {
		logs.PrintCritical(err.Error())
		rr <- true
		return
	}
	defer nfq.Close()
	logs.PrintDebug(fmt.Sprintf(`Opened nfq for service "%v".`, service.Name))

	logs.PrintDebug(fmt.Sprintf(`Starting background task for service "%v"...`, service.Name))
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

		return handlePacket(nfq, &packet, service)
	}, func(err error) int {
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
	logs.PrintDebug(fmt.Sprintf(`Started background task for service "%v".`, service.Name))

	// block this thread until nfq is done handling the packets,
	// basically until forever
	<-ctx.Done()
}

func handlePacket(nfq *nfqueue.Nfqueue, packet *nfqueue.Attribute, service services.Service) int {
	payload := make([]byte, len(*packet.Payload))
	copy(payload, *packet.Payload)

	var offset int
	if service.Protocol == "udp" {
		offset = headers.GetUDPHeaderLength()
	} else if service.Protocol == "tcp" {
		offset = headers.GetTCPHeaderLength(payload)
	}
	payloadString := string(payload)[offset:]

	// join all the regexes with the | operator
	// TODO: check if this actually works properly
	regexMatcher := regexp.MustCompile(strings.Join(service.RegexList, "|"))

	var verdict int
	if regexMatcher.MatchString(payloadString) {
		verdict = nfqueue.NfDrop
	} else {
		verdict = nfqueue.NfAccept
	}
	nfq.SetVerdict(*packet.PacketID, verdict)

	logs.PrintDebug(fmt.Sprintf(`"%v": %v packet %q.`, service.Name, func() string {
		if verdict == nfqueue.NfAccept {
			return "accepted"
		} else if verdict == nfqueue.NfDrop {
			return "dropped"
		}
		return ""
	}(), payloadString))

	// this is a signal to keep receiving messages:
	// https://pkg.go.dev/github.com/florianl/go-nfqueue#ErrorFunc
	return 0
}
