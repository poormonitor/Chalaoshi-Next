package backend

import (
	"fmt"
	"github.com/ccding/go-stun/stun"
	"net"
	"strings"
)

type Addr struct {
	IP   string
	Port uint16
}

// write a addr to string
func (a Addr) String() string {
	if strings.Contains(a.IP, ":") {
		return fmt.Sprintf("[%s]:%d", a.IP, a.Port)
	} else {
		return fmt.Sprintf("%s:%d", a.IP, a.Port)
	}
}

func getAddrStun() (Addr, error) {
	client := stun.NewClient()
	client.SetServerAddr("stun.miwifi.com:3478")

	nat, addr, err := client.Discover()

	if err != nil {
		return Addr{}, err
	}

	if (nat == stun.NATFull) || (nat == stun.NATRestricted) || (nat == stun.NATPortRestricted) {
		return Addr{IP: addr.IP(), Port: addr.Port()}, nil
	}

	return Addr{}, nil
}

func isIPv6Intranet(ip net.IP) bool {
	_, loopbackNet, _ := net.ParseCIDR("::1/128")
	_, linkLocalNet, _ := net.ParseCIDR("fe80::/10")
	_, ulaNet, _ := net.ParseCIDR("fc00::/7")

	if ip.To4() != nil || ip.To16() == nil {
		return false
	}
	return loopbackNet.Contains(ip) || linkLocalNet.Contains(ip) || ulaNet.Contains(ip)
}

func getAddrV6() (Addr, error) {
	const (
		targetIP   = "2400:3200::1" // A public IPv6 DNS server (Alibaba Cloud DNS)
		targetPort = 53             // Standard DNS UDP port
	)

	remoteAddr, err := net.ResolveUDPAddr("udp6", net.JoinHostPort(targetIP, fmt.Sprintf("%d", targetPort)))
	if err != nil {
		return Addr{}, fmt.Errorf("failed to resolve remote address %s:%d: %w", targetIP, targetPort, err)
	}

	conn, err := net.DialUDP("udp6", nil, remoteAddr)
	if err != nil {
		return Addr{}, fmt.Errorf("failed to dial UDP to %s:%d: %w", targetIP, targetPort, err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr()
	udpAddr, ok := localAddr.(*net.UDPAddr)
	if !ok {
		return Addr{}, fmt.Errorf("local address is not a UDP address: %v", localAddr)
	}

	detectedIP := udpAddr.IP

	if detectedIP.To4() != nil || detectedIP.To16() == nil {
		return Addr{}, fmt.Errorf("detected address %s is not a valid IPv6 address or is IPv4-mapped", detectedIP.String())
	}
	// 2. Filter out intranet addresses.
	if isIPv6Intranet(detectedIP) {
		return Addr{}, fmt.Errorf("detected local IPv6 address %s is an intranet address", detectedIP.String())
	}

	hash := uint32(0)
	for i := 0; i < len(detectedIP); i++ {
		hash = hash*31 + uint32(detectedIP[i])
	}
	port := uint16((hash % 30001) + 30000)
	addr := Addr{IP: detectedIP.String(), Port: port}

	return addr, nil
}

func getAddrLocal() []Addr {
	// This function should return the local address in a format suitable for QUIC
	// The addr is a list that can be used, in the order of ipv6 and ipv4
	var addrLocal []Addr
	ch := make(chan Addr, 2)
	go func() {
		addr, err := getAddrStun()
		fmt.Println("Addr STUN:", addr)
		if err == nil && addr.IP != "" {
			ch <- addr
		} else {
			ch <- Addr{}
		}
	}()

	go func() {
		addr, err := getAddrV6()
		fmt.Println("Addr V6:", addr)
		if err == nil && addr.IP != "" {
			ch <- addr
		} else {
			ch <- Addr{}
		}
	}()

	for i := 0; i < 2; i++ {
		addr := <-ch
		if addr.IP != "" {
			addrLocal = append(addrLocal, addr)
		}
	}

	// 如果有两个地址，第二个端口与第一个保持一致
	if len(addrLocal) == 2 {
		addrLocal[1].Port = addrLocal[0].Port
	}

	return addrLocal
}
