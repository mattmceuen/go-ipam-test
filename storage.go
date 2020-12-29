package main

// Based on the default, in-memory storage implementation found here:
// https://github.com/metal-stack/go-ipam/blob/master/memory.go

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	goipam "github.com/metal-stack/go-ipam"
)

type VinoIpamStorage struct {
	prefixes map[string]goipam.Prefix
	lock	 sync.RWMutex
}

// NewMemory create a VinoIpamStorage storage for ipam
func NewVinoIpamStorage() *VinoIpamStorage {
	prefixes := make(map[string]goipam.Prefix)
	return &VinoIpamStorage{
		prefixes: prefixes,
		lock:	 sync.RWMutex{},
	}
}

func (m *VinoIpamStorage) CreatePrefix(prefix goipam.Prefix) (goipam.Prefix, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.prefixes[prefix.Cidr]
	if ok {
		return goipam.Prefix{}, fmt.Errorf("prefix already created:%v", prefix)
	}
	m.prefixes[prefix.Cidr] = *prefix.DeepCopy()
	fmt.Printf("Persisting IPAM to ConfigMap: %v\n", m.prefixes)
	return prefix, nil
}
func (m *VinoIpamStorage) ReadPrefix(prefix string) (goipam.Prefix, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	result, ok := m.prefixes[prefix]
	if !ok {
		return goipam.Prefix{}, errors.Errorf("Prefix %s not found", prefix)
	}
	fmt.Printf("Reading IPAM from ConfigMap: %v\n", m.prefixes)
	return *result.DeepCopy(), nil
}
func (m *VinoIpamStorage) ReadAllPrefixes() ([]goipam.Prefix, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ps := make([]goipam.Prefix, 0, len(m.prefixes))
	for _, v := range m.prefixes {
		ps = append(ps, *v.DeepCopy())
	}
	fmt.Printf("Reading IPAM from ConfigMap: %v\n", m.prefixes)
	return ps, nil
}
func (m *VinoIpamStorage) UpdatePrefix(prefix goipam.Prefix) (goipam.Prefix, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if prefix.Cidr == "" {
		return goipam.Prefix{}, fmt.Errorf("prefix not present:%v", prefix)
	}
	_, ok := m.prefixes[prefix.Cidr]
	if !ok {
		return goipam.Prefix{}, fmt.Errorf("prefix not found:%s", prefix.Cidr)
	}
	m.prefixes[prefix.Cidr] = *prefix.DeepCopy()
	fmt.Printf("Persisting IPAM to ConfigMap: %v\n", m.prefixes)
	return prefix, nil
}
func (m *VinoIpamStorage) DeletePrefix(prefix goipam.Prefix) (goipam.Prefix, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.prefixes, prefix.Cidr)
	fmt.Printf("Persisting IPAM to ConfigMap: %v\n", m.prefixes)
	return *prefix.DeepCopy(), nil
}

