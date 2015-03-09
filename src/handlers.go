// Copyright 2014 go-beacon authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const base64GifPixel = "R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs="

func (s *httpServer) route() {
	http.HandleFunc(s.config.BeaconURI, s.beaconHandler)
	http.HandleFunc("/echo", s.echoBeaconHandler)
	http.Handle("/", http.FileServer(http.Dir(s.config.DocumentRoot)))
}

func (s *httpServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, world\r\n")
}

func (s *httpServer) beaconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-TRACKER-ID", "0")
	w.Header().Set("Content-Type", "image/gif")
	output, _ := base64.StdEncoding.DecodeString(base64GifPixel)
	w.Write(output)
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	if len(r.Form) > 0 {
		go func() {
			err := producer.Send(r.Form)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func (s *httpServer) echoBeaconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("X-TRACKER-ID", "0")
	t, err := json.Marshal(r.URL.Query())
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(t)
}
