// WakeOnLan in Go
// dRbiG, 2014-01-07
// See LICENSE.txt

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"regexp"
	"strings"
)

const (
	CFGNAME = "gowol.json"
)

type Host struct {
	Name, Mac, Broadcast string
}

func mactobyte(mac string) (net.HardwareAddr, error) {
	var validateMac = regexp.MustCompile(`^([[:xdigit:]]{2}:){5}[[:xdigit:]]{2}$`)
	if !validateMac.MatchString(mac) {
		return nil, errors.New("invalid mac-address.")
	}

	return net.ParseMAC(mac)
}

func makepayload(mac string) ([]byte, error) {
	var payload []byte

	bytemac, err := mactobyte(mac)
	if err != nil {
		return nil, err
	}

	payload = make([]byte, 6)
	for i := 0; i < 6; i++ {
		payload[i] = 255
	}

	for i := 0; i < 16; i++ {
		payload = append(payload, bytemac...)
	}

	return payload, nil
}

func wakeup(mac, broadcast string) {
	target, _ := net.ResolveUDPAddr("udp", broadcast+":9")
	sock, _ := net.DialUDP("udp", nil, target)

	payload, err := makepayload(mac)
	if err != nil {
		log.Fatal(err)
	}

	sock.Write(payload)
}

func main() {
	usr, _ := user.Current()
	cfgfilename := usr.HomeDir + "/" + CFGNAME
	cfgfile, err := os.Open(cfgfilename)
	if os.IsNotExist(err) {
		// create the file and open it.
		cfgfile, err = os.Create(cfgfilename)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}
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
