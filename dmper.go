package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Hold the data
var mhostport = map[string]int{}

// Find default network interface
// No functions in go. Use netstat or parse proc
func get_device() (string, error) {
	path, err := exec.LookPath("netstat")
	if err != nil {
		return "", err
	}
	cmd, err := exec.Command(path, "-nr").Output()
	if err != nil {
		return "", err
	}

	var device string
	for _, i := range strings.Split(string(cmd), "\n") {
		fields := strings.Fields(i)
		if fields[0] == "0.0.0.0" {
			device = fields[len(fields)-1]
			break
		}
	}
	if device == "" {

		return "", fmt.Errorf("No devices found")
	}
	return device, nil

}

// An http server to return stats
func runhttpserver(port, readtimeout, writetimeout int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/dumpmetrics", dumphandler)
	mux.HandleFunc("/dumpandclear", clearhandler)
	svr := &http.Server{
		Addr:           ":" + strconv.Itoa(port),
		ReadTimeout:    time.Duration(readtimeout) * time.Second,
		WriteTimeout:   time.Duration(writetimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}
	log.Fatal(svr.ListenAndServe())
}

// Dump a json of metrics collected
func dumphandler(w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(mhostport)
	if err != nil {
		//log.Printf("Got error: ", err)
		fmt.Fprintf(w, "ERROR: %s", err.Error())
	}
	fmt.Fprintf(w, string(js))
}

// Dump a json of metrics collected, clear it
func clearhandler(w http.ResponseWriter, r *http.Request) {
	dumphandler(w, r)
	mhostport = map[string]int{}

}

func main() {
	var port, writetimeout, readtimeout int
	var verbose bool
	flag.IntVar(&port, "port", 8080, "Port for http")
	flag.IntVar(&writetimeout, "writetimeout", 10, "Write timeout for http")
	flag.IntVar(&readtimeout, "readtimeout", 10, "Read timeout for htt&p")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode")
	flag.Parse()

	dev, err := get_device()
	if err != nil {
		log.Fatal("Unable to get PCAP device: ", err)
	}
	if verbose == true {
		fmt.Println("Got dev: ", dev)
	}
	go runhttpserver(port, readtimeout, writetimeout)

	// Get hostname. This is needed for filter
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Unable to get hostname: ", err)
	}
	// PCAP filter for dump
	var filter string = "tcp and dst host not " + hostname

	handle, err := pcap.OpenLive(dev, 0, true, 5*time.Second)
	if err != nil {
		log.Fatal("Pcap failed: ", err)
	}
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Println("Unable to set filter: ", err)

	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		mhostport = setpacketinfo(packet, mhostport)
		if verbose == true {
			fmt.Printf("%#v\n", len(mhostport))
		}
	}
}

// Sets info
func setpacketinfo(packet gopacket.Packet, m map[string]int) map[string]int {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	var host string
	var port int
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		//fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		host = fmt.Sprintf("%s", ip.DstIP)
	}
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		//fmt.Printf("From %d to %d\n", tcp.SrcPort, tcp.DstPort)
		p := fmt.Sprintf("%d", tcp.DstPort)
		port, _ = strconv.Atoi(p)
	}
	m[host] = port

	return m
}
