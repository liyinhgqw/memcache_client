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

package memcache

import (
	"fmt"
	"hash/crc32"
	"net"
	"sort"
	"strings"
	"sync"
)

const (
	maxFactor = 50
)

type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// ChServerList is consistent hashing ServerSelector.
type ChServerList struct {
	lk   sync.RWMutex
	vn   []uint32
	addr map[uint32]net.Addr
}

type WeightedServer struct {
	Addr   string
	Weight uint32
}

// SetServers changes a ChServerList's set of servers at runtime and is
// threadsafe.
//
// SetServers returns an error if any of the server names fail to
// resolve. No attempt is made to connect to the server. If any error
// is returned, no changes are made to the ServerList.
func (ss *ChServerList) SetServers(servers ...*WeightedServer) error {
	var totalWeight uint32 = 0
	totalServers := len(servers)
	nhash := make(map[uint32]net.Addr)

	for _, server := range servers {
		if server.Weight == 0 {
			server.Weight++
		}

		totalWeight += server.Weight

	}

	for _, server := range servers {
		var (
			addr net.Addr
			err  error
		)
		if strings.Contains(server.Addr, "/") {
			addr, err = net.ResolveUnixAddr("unix", server.Addr)
			if err != nil {
				return err
			}
		} else {
			addr, err = net.ResolveTCPAddr("tcp", server.Addr)
			if err != nil {
				return err
			}
		}
		// compute how many virtual nodes for this server
		vnc := int(float32(maxFactor*uint32(totalServers)*server.Weight) / float32(totalWeight))
		for i := 0; i < vnc; i++ {
			vnname := fmt.Sprintf("%s%d", server.Addr, i)
			h := ss.hash(vnname)
			nhash[h] = addr
		}
	}

	var nvn []uint32
	for i := range nhash {
		nvn = append(nvn, i)
	}
	sort.Sort(Uint32Slice(nvn))

	ss.lk.Lock()
	defer ss.lk.Unlock()
	ss.vn = nvn
	ss.addr = nhash

	return nil
}

func (ss *ChServerList) PickServer(key string) (net.Addr, error) {
	ss.lk.RLock()
	defer ss.lk.RUnlock()
	if len(ss.addr) == 0 || len(ss.vn) == 0 {
		return nil, ErrNoServers
	}
	cs := ss.hash(key)
	serv := sort.Search(len(ss.addr), func(i int) bool { return ss.vn[i] >= cs })
	if serv >= len(ss.addr) {
		serv = 0
	}
	return ss.addr[ss.vn[serv]], nil
}

func (ss *ChServerList) hash(key string) uint32 {
	cs := crc32.ChecksumIEEE([]byte(key))
	return cs
}
