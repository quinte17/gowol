// WakeOnLan in Go
// dRbiG, 2014-01-07
// See LICENSE.txt

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"log"
	"io"
	"strings"
)

const (
	CFGNAME = "gowol.json"
)

type Host struct {
	Name, Mac, Broadcast string
}

var (
	payload []byte
)

func makepayload(mac string) bool {
	var x byte

	if len(mac) != 17 {
		return false
	}

	for i := 0; i < 6; i++ {
		_, err := fmt.Sscanf(mac[3*i:3*i+2], "%2X", &x)
		if err != nil {
			return false
		}
		for c := 0; c < 16; c++ {
			payload[6*c+6+i] = x
		}
	}

	return true
}

func init() {
	// initialize payload
	payload = make([]byte, 102)
	for i := 0; i < 6; i++ {
		payload[i] = 255
	}
}

func wakeup(mac, broadcast string) {
	target, _ := net.ResolveUDPAddr("udp", broadcast+":9")
	sock, _ := net.DialUDP("udp", nil, target)

	if !makepayload(mac) {
		log.Fatal("mac parse error:", mac)
	}
	sock.Write(payload)
}

func main() {
	cfgfile, _ := os.Open(CFGNAME)
	dec := json.NewDecoder(cfgfile)
	for {
		var host Host
		if err := dec.Decode(&host); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if strings.ToLower(host.Name) == strings.ToLower(os.Args[1]) {
			fmt.Println("waking up", host.Name, host.Mac, host.Broadcast)
			wakeup(host.Mac, host.Broadcast)
		}
	}
}
