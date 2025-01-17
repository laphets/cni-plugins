// Copyright 2015 CNI authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing

import (
	"net"
	"os"

	"github.com/containernetworking/plugins/plugins/ipam/host-local/backend"
)

type FakeStore struct {
	ipMap          map[string]string
	lastReservedIP map[string]net.IP
}

// FakeStore implements the Store interface
var _ backend.Store = &FakeStore{}

func NewFakeStore(ipmap map[string]string, lastIPs map[string]net.IP) *FakeStore {
	return &FakeStore{ipmap, lastIPs}
}

func (s *FakeStore) Lock() error {
	return nil
}

func (s *FakeStore) Unlock() error {
	return nil
}

func (s *FakeStore) Close() error {
	return nil
}

func (s *FakeStore) Reserve(id string, ifname string, ip net.IP, rangeID string) (bool, error) {
	key := ip.String()
	if _, ok := s.ipMap[key]; !ok {
		s.ipMap[key] = id
		s.lastReservedIP[rangeID] = ip
		return true, nil
	}
	return false, nil
}

func (s *FakeStore) LastReservedIP(rangeID string) (net.IP, error) {
	ip, ok := s.lastReservedIP[rangeID]
	if !ok {
		return nil, os.ErrNotExist
	}
	return ip, nil
}

func (s *FakeStore) Release(ip net.IP) error {
	delete(s.ipMap, ip.String())
	return nil
}

func (s *FakeStore) ReleaseByID(id string, ifname string) error {
	toDelete := []string{}
	for k, v := range s.ipMap {
		if v == id {
			toDelete = append(toDelete, k)
		}
	}
	for _, ip := range toDelete {
		delete(s.ipMap, ip)
	}
	return nil
}

func (s *FakeStore) GetByID(id string, ifname string) []net.IP {
	var ips []net.IP
	for k, v := range s.ipMap {
		if v == id {
			ips = append(ips, net.ParseIP(k))
		}
	}
	return ips
}

func (s *FakeStore) SetIPMap(m map[string]string) {
	s.ipMap = m
}

func (s *FakeStore) RecordRelease(ip net.IP, rangeID string) error {
	return nil
}
func (s *FakeStore) LastReleasedIP(rangeID string) (net.IP, error) {
	return nil, nil
}
