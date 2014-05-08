/*
Copyright 2011 Google Inc and Yin Li.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package memcache provides a client for the memcached cache server.
package memcache

import (
	"net"
	"testing"
)

var testServers = []string{"localhost:11211", "localhost:11212"}

func setupCh(t *testing.T) bool {
	for _, ts := range testServers {
		c, err := net.Dial("tcp", ts)
		if err != nil {
			t.Logf("skipping test; no server running at %s", testServer)
			return false
		}
		c.Write([]byte("flush_all\r\n"))
		c.Close()
	}
	return true
}

func TestChLocalhost(t *testing.T) {
	if !setupCh(t) {
		return
	}
	testWithClient(t, NewChClient(testServers, nil))
}
