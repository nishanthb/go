package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	/*
	 */)

// Hold the data
var (
	mhostport = map[string]map[int]int{}
	mapLock   sync.Mutex
)

// Find default network interface
// No functions in go. Use netstat or parse proc
func getdevice() (string, error) {
	data, err := ioutil.ReadFile("/proc/net/route")
	if err != nil {
		return "", err
	}
	var device string = ""
	// We slacked. We should check for 0003 and take the lowest Metric field
	// for _,i := range strings.Split(string(data),"\n") {
	//	fields := strings.Fields(i)
	//	if fields[3] == "0003" {
	//
	//		device = fields[0]
	//		break
	//	}
	//}
	mroute := make(map[string]int)
	for _, i := range strings.Split(string(data), "\n") {
		if len(i) == 0 {
			continue
		}
		fields := strings.Fields(i)
		if fields[3] == "0003" {
			num, err := strconv.Atoi(fields[6])
			if err != nil {
				continue
			}
			mroute[fields[0]] = num
		}
	}
	device = getsortedmin(mroute)
	if device == "" {
		return "", fmt.Errorf("Unable to determine gateway device")
	}
	return device, nil
}

// sort map
func getsortedmin(r map[string]int) string {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range r {
		sorted = append(sorted, kv{k, v})
	}
	// Requires go 1.8
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	return sorted[len(sorted)-1].Key

}

// An http server to return stats
func runhttpserver(port, readtimeout, writetimeout int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/dump", dumphandler)
	mux.HandleFunc("/clear", clearhandler)
	svr := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  time.Duration(readtimeout) * time.Second,
		WriteTimeout: time.Duration(writetimeout) * time.Second,
		IdleTimeout:  time.Duration(10) * time.Second,
		Handler:      mux,
	}
	log.Fatal(svr.ListenAndServe())
}

// Dump a json of metrics collected
func dumphandler(w http.ResponseWriter, r *http.Request) {
	mapLock.Lock()
	js, err := json.Marshal(mhostport)
	mapLock.Unlock()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Fprintf(w, string(js))
}

// Dump a json of metrics collected, clear it
func clearhandler(w http.ResponseWriter, r *http.Request) {
	dumphandler(w, r)
	//mapLock.Lock()
	mhostport = map[string]map[int]int{}
	//mapLock.Unlock()
}

func main() {

	var port, writetimeout, readtimeout int
	var verbose bool
	flag.IntVar(&port, "port", 8080, "Port for http")
	flag.IntVar(&writetimeout, "writetimeout", 10, "Write timeout for http")
	flag.IntVar(&readtimeout, "readtimeout", 10, "Read timeout for htt&p")
	//flag.IntVar(&progtimeout, "progtimeout", 300, "Exit after this much time")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode")
	flag.Parse()

	if verbose == true {
		fmt.Println("Our PID: ", strconv.Itoa(os.Getpid()))
	}
	dev, err := getdevice()
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

	// Whats an ideal snaplen?
	handle, err := pcap.OpenLive(dev, 1024, true, 5*time.Second)
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
	}
}

// Sets info
func setpacketinfo(packet gopacket.Packet, m map[string]map[int]int) map[string]map[int]int {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	var host string
	var port int
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		host = fmt.Sprintf("%s", ip.DstIP)
	}
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		p := fmt.Sprintf("%d", tcp.DstPort)
		port, _ = strconv.Atoi(p)
	}
	m = add(m, host, port)

	return m
}

// Add to multi dimensional map, copied off blog.golang.org/go-maps-in-action
// Though they recommend struct usage
func add(m map[string]map[int]int, host string, port int) map[string]map[int]int {
	//mapLock.Lock()
	mm, ok := m[host]
	if !ok {
		mm = make(map[int]int)
		m[host] = mm
	}
	mm[port] = 1
	//mm[port]++
	//mapLock.Unlock()
	return m
}
