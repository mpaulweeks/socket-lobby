// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
)

var addr = flag.String("addr", "localhost:5110", "http service address")

func main() {
	flag.Parse()

	h := newHandler()
	s := newServer(*addr, h)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
