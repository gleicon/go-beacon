package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fiorix/go-redis/redis"
	"github.com/gdamore/mangos"
	"github.com/gdamore/mangos/protocol/rep"
	"github.com/gdamore/mangos/transport/ipc"
	"github.com/gdamore/mangos/transport/tcp"
	"github.com/ugorji/go/codec"
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

func boomerangMetrics(rc *redis.Client, prefix string, d map[string][]string, pbkey string) {
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

	doc := make(map[string]interface{})
	doc[fmt.Sprintf("%s.navigation.type", partial)] = ntNavtype
	doc[fmt.Sprintf("%s.navigation.timing.dns", partial)] = ntDNS
	doc[fmt.Sprintf("%s.navigation.timing.connection", partial)] = ntCon
	doc[fmt.Sprintf("%s.navigation.timing.dom.loaded", partial)] = ntDomcontloaded
	doc[fmt.Sprintf("%s.navigation.timing.dom.processing", partial)] = ntProcessed
	doc[fmt.Sprintf("%s.navigation.timing.request", partial)] = ntRequest
	doc[fmt.Sprintf("%s.navigation.timing.response", partial)] = ntResponse
	doc[fmt.Sprintf("%s.roundtrip", partial)] = roundtrip
	doc[fmt.Sprintf("%s.page", partial)] = page

	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}
	fmt.Println(rc.RPush(pbkey, string(jsonDoc)))
}

func jsMetrics(rc *redis.Client, prefix string, d map[string][]string, pbkey string) {
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

func listenMangos(listenAddr *string, trackerType *string, prefix *string, rc *redis.Client, pbkey string) {
	var err error
	var responseServer mangos.Socket
	responseServer, err = rep.NewSocket()
	responseServer.AddTransport(ipc.NewTransport())
	responseServer.AddTransport(tcp.NewTransport())

	if err = responseServer.Listen(*listenAddr); err != nil {
		fmt.Printf("\nServer listen failed: %v", err)
		return
	}

	fmt.Println("Listening")
	mh.MapType = reflect.TypeOf(map[string][]string(nil))
	go func() {
		for {
			serverMsg, err := responseServer.RecvMsg()
			if err != nil {
				fmt.Printf("\nServer receive failed: %v", err)
			}
			d, err := decode(serverMsg.Body)
			if len(d) < 1 {
				fmt.Println("Discarded message")
				continue
			}
			switch *trackerType {
			case "boomerang":
				boomerangMetrics(rc, *prefix, d, pbkey)
			case "js":
				jsMetrics(rc, *prefix, d, pbkey)
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
	// pbcollector -type boomerang -listen tcp://127.0.0.1:8000 -redis 192.168.33.20:8125

	listenAddr := flag.String("listen", "tcp://127.0.0.1:8000", "Listening string - default: tcp://127.0.0.1:8000")
	redisServer := flag.String("packetbeat", "127.0.0.1:6380", "packet beat redis - default: 127.0.0.1:6380")
	trackerType := flag.String("tracker", "boomerang", "tracker type - default: boomerang [boomerang, js]")
	redisPBKey := flag.String("redispbkey", "packetbeat", "packetbeat redis' key - default: packetbeat")
	prefix := flag.String("prefix", "", "prefix to metrics - default: empty string")
	flag.Parse()

	fmt.Println("Listening:", *listenAddr)
	fmt.Println("Redis endpoint:", *redisServer)
	fmt.Println("Tracker type: ", *trackerType)
	fmt.Println("Prefix: ", *prefix)
	fmt.Println("Packetbeat key: ", *redisPBKey)
	fmt.Println("Consumer ready")
	rc := redis.New(*redisServer)

	listenMangos(listenAddr, trackerType, prefix, rc, *redisPBKey)

	// wait for cleanup
	for {
		time.Sleep(10 * time.Second)
	}
}
