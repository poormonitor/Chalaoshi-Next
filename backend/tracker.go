package backend

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/xgfone/go-bt/metainfo"
	"github.com/xgfone/go-bt/tracker"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	getTransport = func(myAddr string) *http.Transport {
		return &http.Transport{
			DialContext: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				var networkNew = network
				if strings.Contains(myAddr, ":") {
					networkNew += "6"
				} else {
					networkNew += "4"
				}

				var conn net.Conn
				if strings.HasPrefix(network, "udp") {
					udpAddr, err := net.ResolveUDPAddr(networkNew, addr)
					if err != nil {
						return nil, err
					}
					conn, err = net.DialUDP(networkNew, nil, udpAddr)
					if err != nil {
						return nil, err
					}
				} else {
					tcpAddr, err := net.ResolveTCPAddr(networkNew, addr)
					if err != nil {
						return nil, err
					}
					conn, err = net.DialTCP(networkNew, nil, tcpAddr)
					if err != nil {
						return nil, err
					}
				}

				return conn, nil
			}}
	}
)

func getTrackerList() []string {
	var rawText string
	err := requests.URL("https://cf.trackerslist.com/best.txt").ToString(&rawText).Fetch(context.Background())
	if err != nil {
		return nil
	}
	lines := strings.Split(rawText, "\n")
	return lines
}

func requestTracker(trackerAddr string, addr Addr, peerId metainfo.Hash, infoHash metainfo.Hash) {
	httpClient := http.Client{Transport: getTransport(addr.IP)}
	client, err := tracker.NewClient(trackerAddr, metainfo.NewRandomHash(), &httpClient)
	if err != nil {
		return
	}

	request := tracker.AnnounceRequest{IP: net.ParseIP(addr.IP), Port: addr.Port, Downloaded: 1, PeerID: peerId, InfoHash: infoHash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = client.Announce(ctx, request)
	if err != nil {
		return
	}
	fmt.Printf("Announced to tracker %s with addr %s\n", trackerAddr, addr.String())
}

func announcePeer(addrLocal []Addr) {
	trackers := getTrackerList()
	peerId := metainfo.NewRandomHash()
	infoHash := metainfo.NewHashFromString("chalaoshi-next")

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		for _, trackerAddr := range trackers {
			if !strings.HasPrefix(trackerAddr, "udp") && !strings.HasPrefix(trackerAddr, "http") {
				continue
			}

			for _, addr := range addrLocal {
				go requestTracker(trackerAddr, addr, peerId, infoHash)
			}
		}
		<-ticker.C
	}
}

func requestPeer(trackerAddr string, peerId metainfo.Hash, infoHash metainfo.Hash) ([]Addr, error) {
	httpClient := http.Client{Transport: getTransport("127.0.0.1")}
	client, err := tracker.NewClient(trackerAddr, metainfo.NewRandomHash(), &httpClient)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	request := tracker.AnnounceRequest{IP: net.ParseIP("::1"), Port: 1234, Downloaded: 0, PeerID: peerId, InfoHash: infoHash}
	resp, err := client.Announce(ctx, request)

	if err != nil {
		return nil, err
	}

	if resp.Addresses != nil {
		var peersAddr []Addr
		for _, peer := range resp.Addresses {
			peersAddr = append(peersAddr, Addr{IP: peer.Host, Port: peer.Port})
		}
		return peersAddr, nil
	} else {
		return nil, nil
	}
}

func getPeer() (addrLocal []Addr) {
	trackers := getTrackerList()
	peerId := metainfo.NewRandomHash()
	infoHash := metainfo.NewHashFromString("chalaoshi-next")

	type result struct {
		peers []Addr
	}
	resultCh := make(chan result, len(trackers))

	trackerCnt := 0
	for _, trackerAddr := range trackers {
		if !strings.HasPrefix(trackerAddr, "udp") && !strings.HasPrefix(trackerAddr, "http") {
			continue
		}
		trackerAddrCopy := trackerAddr
		trackerCnt++
		go func() {
			peersAddr, err := requestPeer(trackerAddrCopy, peerId, infoHash)
			if err == nil && len(peersAddr) > 0 {
				resultCh <- result{peers: peersAddr}
			} else {
				resultCh <- result{peers: nil}
			}
		}()
	}

	for i := 0; i < trackerCnt; i++ {
		res := <-resultCh
		if len(res.peers) > 0 {
			addrLocal = append(addrLocal, res.peers...)
		}
	}

	// 去重
	uniquePeers := make(map[string]Addr)
	for _, peer := range addrLocal {
		key := peer.String()
		if _, exists := uniquePeers[key]; !exists && peer.Port != 1234 {
			uniquePeers[key] = peer
		}
	}

	var resultList []Addr
	for _, peer := range uniquePeers {
		resultList = append(resultList, peer)
	}
	return resultList
}
