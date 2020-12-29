package main

import (
	"fmt"
	goipam "github.com/metal-stack/go-ipam"
)

func main() {
	// VINO needs to do the following things:

	// 1. persist IPAM status, e.g. to a configmap or CR
	fmt.Println("** Initializing IPAM from pretend ConfigMap storage")
	storage := NewVinoIpamStorage()
	ipam := goipam.NewWithStorage(storage)

	// 2. be configured with an IPAM range
	fmt.Println("** Configuring IPAM range")
	prefix, err := ipam.NewPrefix("192.168.0.0/24")
	if err != nil {
		panic(err)
	}

	// 3. allocate IPs when creating BMH cloud-init contents
	fmt.Println("** Acquiring IP")
	ip, err := ipam.AcquireIP(prefix.Cidr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("acquired IP: %s\n", ip.IP)

	// 4. deallocate IPs when deprovisioning vBMH
	fmt.Println("** Releasing IP")
	prefix, err = ipam.ReleaseIP(ip)
	if err != nil {
		panic(err)
	}
	fmt.Printf("released IP: %s.\n", ip.IP)
}

