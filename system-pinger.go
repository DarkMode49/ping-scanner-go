package main

import (
	"time"

	"github.com/prometheus-community/pro-bing"
)

type SystemPinger struct {
	Timeout time.Duration
}

func (p *SystemPinger) Ping(ipAddress string) bool {
	pinger, pingerError := probing.NewPinger(ipAddress)
	if pingerError != nil {
		logMsg("%v", pingerError)
	}
	pinger.Count = 1
	pinger.Timeout = p.Timeout
	pinger.SetPrivileged(true) //? Super-user requirement

	pinger.OnRecv = func(pkt *probing.Packet) {
		logMsg("%d bytes from %s: icmp_seq=%d time=%v",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pingerError = pinger.Run()
	if pingerError != nil {
		logMsg("%v", pingerError)
	}

	return pinger.PacketsRecv > 0
}
