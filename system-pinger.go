package main

import (
	"time"

	"github.com/go-ping/ping"
)

type SystemPinger struct {
	Timeout time.Duration
}

func (p *SystemPinger) Ping(ipAddress string) (PingResult, error) {
	pingResult := PingResult{}
	
	pinger, pingerError := ping.NewPinger(ipAddress)
	if pingerError != nil {
		return PingResult{}, pingerError 
	}
	pinger.Count = 1
	pinger.Timeout = p.Timeout
	pinger.SetPrivileged(true) //? Super-user requirement

	pinger.OnRecv = func(pkt *ping.Packet) {
		pingResult.Bytes = pkt.Nbytes
		pingResult.IPAddr = pkt.IPAddr
		pingResult.Sequence = pkt.Seq
		pingResult.Latency = pkt.Rtt
	}

	pingerError = pinger.Run()
	
	if pingerError != nil {
		return PingResult{}, pingerError
	}

	return pingResult, nil
}
