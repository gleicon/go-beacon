package main

import (
	"bitbucket.org/gdamore/mangos"
	"bitbucket.org/gdamore/mangos/protocol/rep"
	"bitbucket.org/gdamore/mangos/transport/all"
	"flag"
	"fmt"
	"github.com/ugorji/go/codec"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	mh codec.MsgpackHandle
	b  []byte
)

func boomerangMetrics(udpAddr *net.UDPAddr, prefix string, d map[string][]string) {
	fmt.Println("------ server msg ------ ")
	ntDNS, _ := delta(d["nt_dns_st"][0], d["nt_dns_end"][0])                               // domainLookupEnd - domainLookupStart
	ntCon, _ := delta(d["nt_con_st"][0], d["nt_con_end"][0])                               // connectEnd - connectStart
	ntDomcontloaded, _ := delta(d["nt_domcontloaded_st"][0], d["nt_domcontloaded_end"][0]) // domContentLoadedEnd - domContentLoadedStart
	ntProcessed, _ := delta(d["nt_domcontloaded_st"][0], d["nt_domcomp"][0])               // domComplete - domContentLoadedStart
	ntRequest, _ := delta(d["nt_req_st"][0], d["nt_res_st"][0])                            // ResponseStart - RequestStart
	ntResponse, _ := delta(d["nt_res_st"][0], d["nt_res_end"][0])                          // ResponseEnd - ResponseStart
	ntNavtype := d["nt_nav_type"][0]
	roundtrip, _ := delta(d["rt.bstart"][0], d["rt.end"][0])
	page := d["r"][0]
	url, err := url.Parse(d["u"][0])
	if err != nil {
		fmt.Println("Error parsing URL", err)
		return
	}
	if prefix != "" {
		prefix = prefix + "."
	}
	partial := fmt.Sprintf("%s%s%s", prefix, strings.Replace(url.Host, "/", ".", -1), strings.Replace(url.Path, "/", ".", -1))
	partial = strings.TrimSuffix(partial, ".")

	fmt.Printf("%s.navigation.type: %s\n", partial, ntNavtype)
	fmt.Printf("%s.navigation.timing.dns: %d\n", partial, ntDNS)
	fmt.Printf("%s.navigation.timing.connection: %d\n", partial, ntCon)
	fmt.Printf("%s.navigation.timing.dom.loaded: %d\n", partial, ntDomcontloaded)
	fmt.Printf("%s.navigation.timing.dom.processing: %d\n", partial, ntProcessed)
	fmt.Printf("%s.navigation.timing.request: %d\n", partial, ntRequest)
	fmt.Printf("%s.navigation.timing.response: %d\n", partial, ntResponse)
	fmt.Printf("%s.roundtrip: %d\n", partial, roundtrip)
	fmt.Printf("%s.page: %s\n", partial, page)
	fmt.Println("------ server msg ------ ")
}

func jsMetrics(udpAddr *net.UDPAddr, prefix string, d map[string][]string) {
	fmt.Println("------ server msg ------ ")
	ntDNS, _ := delta(d["nt_dns_st"][0], d["nt_dns_end"][0])                               // domainLookupEnd - domainLookupStart
	ntCon, _ := delta(d["nt_con_st"][0], d["nt_con_end"][0])                               // connectEnd - connectStart
	ntDomcontloaded, _ := delta(d["nt_domcontloaded_st"][0], d["nt_domcontloaded_end"][0]) // domContentLoadedEnd - domContentLoadedStart
	ntProcessed, _ := delta(d["nt_domcontloaded_st"][0], d["nt_domcomp"][0])               // domComplete - domContentLoadedStart
	ntRequest, _ := delta(d["nt_req_st"][0], d["nt_res_st"][0])                            // ResponseStart - RequestStart
	ntResponse, _ := delta(d["nt_res_st"][0], d["nt_res_end"][0])                          // ResponseEnd - ResponseStart
	ntNavtype := d["nt_nav_type"][0]
	roundtrip, _ := delta(d["rt.bstart"][0], d["rt.end"][0])
	page := d["r"][0]
	url := d["u"][0]

	fmt.Println("Navigation type: ", ntNavtype)
	fmt.Println("Navigation timing DNS: ", ntDNS)
	fmt.Println("Navigation timing Connection: ", ntCon)
	fmt.Println("Navigation timing DOM content loaded: ", ntDomcontloaded)
	fmt.Println("Navigation timing DOM processing: ", ntProcessed)
	fmt.Println("Navigation timing Request: ", ntRequest)
	fmt.Println("Navigation timing Response: ", ntResponse)
	fmt.Println("Roundtrip: ", roundtrip)
	fmt.Println("Page: ", page)
	fmt.Println("URL: ", url)
	fmt.Println("------ server msg ------ ")

}

// Calculate delta between start and end
func delta(start string, end string) (int, error) {
	s, err := strconv.Atoi(start)
	if err != nil {
		return -1, err
	}
	e, err := strconv.Atoi(end)
	if err != nil {
		return -1, err
	}
	return e - s, nil
}

func decode(buf []byte) (map[string][]string, error) {

	doc := map[string][]string(nil)
	dec := codec.NewDecoderBytes(buf, &mh)
	err := dec.Decode(&doc)

	if err != nil {
		return nil, err
	}
	return doc, nil
}

func listenMangos(listenAddr *string, trackerType *string, udpAddr *net.UDPAddr, prefix *string) {
	var err error
	responseServerReady := make(chan struct{})
	responseServer, err := rep.NewSocket()
	defer responseServer.Close()
	var serverMsg *mangos.Message

	if err = responseServer.Listen(*listenAddr); err != nil {
		fmt.Printf("\nServer listen failed: %v", err)
		return
	}

	close(responseServerReady)

	all.AddTransports(responseServer)
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}

	fmt.Println("Listening")
	mh.MapType = reflect.TypeOf(map[string][]string(nil))
	go func() {
		for {
			if serverMsg, err = responseServer.RecvMsg(); err != nil {
				fmt.Printf("\nServer receive failed: %v", err)
			}
			d, err := decode(serverMsg.Body)
			if len(d) < 1 {
				fmt.Println("Discarded message")
				continue
			}
			switch *trackerType {
			case "boomerang":
				boomerangMetrics(udpAddr, *prefix, d)
			case "js":
				jsMetrics(udpAddr, *prefix, d)
			}

			serverMsg.Body = []byte("OK")
			if err = responseServer.SendMsg(serverMsg); err != nil {
				fmt.Printf("\nServer send failed: %v", err)
				return
			}
		}
	}()
}

func main() {
	// consumer -type boomerang -listen tcp://127.0.0.1:8000 -statsd 192.168.33.20:8125

	listenAddr := flag.String("listen", "tcp://127.0.0.1:8000", "Listening string - default: tcp://127.0.0.1:8000")
	statsdServer := flag.String("statsd", "127.0.0.1:8125", "statsd endpoint - default: 127.0.0.1:8125")
	trackerType := flag.String("tracker", "boomerang", "tracker type - default: boomerang [boomerang, js]")
	prefix := flag.String("prefix", "", "prefix to metrics - default: empty string")
	flag.Parse()

	fmt.Println("Listening:", *listenAddr)
	fmt.Println("Statsd endpoint:", *statsdServer)
	fmt.Println("Tracker type: ", *trackerType)
	fmt.Println("Prefix: ", *prefix)
	fmt.Println("Consumer ready")
	udpAddr, err := net.ResolveUDPAddr("udp4", *statsdServer)

	if err != nil {
		fmt.Println("Error resolving statsd server", err)
		return
	}

	listenMangos(listenAddr, trackerType, udpAddr, prefix)

	// wait for cleanup
	for {
		time.Sleep(10 * time.Second)
	}
}
