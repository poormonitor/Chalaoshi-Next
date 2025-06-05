package backend

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"net"
	"strings"
)

var chunkSize = 1024 * 1024 // 1 MB

func runServer(port uint16) {
	c := context.Background()

	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{Port: int(port)})
	if err != nil {
		fmt.Println("Failed to listen on UDP port:", err)
		return
	}

	tr := quic.Transport{Conn: udpConn}
	tlsConf, err := certSetup()
	if err != nil {
		fmt.Println("Failed to set up TLS configuration:", err)
		return
	}

	quicConf := quic.Config{}
	ln, err := tr.Listen(tlsConf, &quicConf)
	if err != nil {
		fmt.Println("Failed to start QUIC listener:", err)
		return
	}

	for {
		conn, err := ln.Accept(c)
		if err == nil {
			go handleConn(conn)
		}
	}
}

func handleConn(conn quic.Connection) {
	// This function can be used to handle the connection if needed
	c := context.Background()
	fmt.Println("New connection from:", conn.RemoteAddr())
	stream, err := conn.AcceptStream(c)
	if err != nil {
		fmt.Println("Failed to accept stream:", err)
		return
	}

	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if n == 0 && err != nil {
		fmt.Println("Failed to read from stream:", err)
		return
	}

	str := string(buf[:n])
	fmt.Println("Content:", str)
	stream, err = conn.OpenStream()
	if err != nil {
		fmt.Println("Failed to open stream:", err)
		return
	}

	defer stream.Close()

	if status != 0 {
		stream.Write([]byte("forbidden"))
		return
	}

	if strings.HasPrefix(str, "hel") {
		if strings.HasSuffix(str, "chalaoshi-next-client") {
			if connected {
				stream.Write([]byte("hel chalaoshi-next-server-full"))
			} else {
				stream.Write([]byte("hel chalaoshi-next-server"))
			}
		} else {
			stream.Write([]byte("forbidden"))
		}
	} else if strings.HasPrefix(str, "dow") {
		fileBytes := getFileBytes(filePath)
		for offset := 0; offset < len(fileBytes); offset += chunkSize {
			end := offset + chunkSize
			if end > len(fileBytes) {
				end = len(fileBytes)
			}
			_, err = stream.Write(fileBytes[offset:end])
			if err != nil {
				fmt.Println("Failed to write chunk to stream:", err)
				return
			}
		}
	} else if strings.HasPrefix(str, "req") {
		if !connected {
			stream.Write([]byte("forbidden"))
			return
		}
		// split the string get the content after the first space
		parts := strings.SplitN(str, " ", 2)
		if len(parts) < 2 {
			stream.Write([]byte("forbidden"))
		} else {
			content := parts[1]
			stream.Write([]byte("success " + content))
		}
	}

}

func handleRequest(request string, conn quic.Connection) string {
	stream, err := conn.OpenStream()
	if err != nil {
		fmt.Println("Failed to open stream:", err)
		return ""
	}

	n, err := stream.Write([]byte(request))
	if err != nil {
		fmt.Println("Failed to write to stream:", err)
		return ""
	}
	stream.Close()

	stream, err = conn.AcceptStream(context.Background())
	if err != nil {
		fmt.Println("Failed to accept stream:", err)
		return ""
	}

	buf := make([]byte, 1024)
	n, err = stream.Read(buf)
	if n == 0 && err != nil {
		fmt.Println("Failed to read from stream:", err)
		return ""
	}

	fmt.Println(n, "bytes read from stream")
	response := string(buf[:n])
	fmt.Println("Received response:", response)

	return response
}

func downloadFileFromPeer() error {
	var peer *Addr

	status = 3
	peers := getPeer()
	fmt.Println("Found peers:", peers)

	status = 4
	peer = tryPeer(peers, false)

	if peer == nil {
		fmt.Println("No suitable peer found")
		return fmt.Errorf("no suitable peer found")
	}

	status = 5
	ctx := context.Background()
	udpAddr, err := net.ResolveUDPAddr("udp", peer.String())
	if err != nil {
		return fmt.Errorf("failed to resolve address: %w", err)
	}

	conn, err := quic.DialAddr(ctx, udpAddr.String(), &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	stream, err := conn.OpenStream()
	if err != nil {
		return fmt.Errorf("failed to open stream: %w", err)
	}

	_, err = stream.Write([]byte("dow"))
	if err != nil {
		return fmt.Errorf("failed to write to stream: %w", err)
	}
	stream.Close()
	stream, err = conn.AcceptStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to accept stream: %w", err)
	}

	fileData := make([]byte, 0, chunkSize*100)
	buf := make([]byte, chunkSize*2)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			fileData = append(fileData, buf[:n]...)
		}
		if n == 0 && err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read from stream: %w", err)
		}
	}

	zipFile, err := createZipFile()
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	if _, err = zipFile.Write(fileData); err != nil {
		return fmt.Errorf("failed to write downloaded content to zip file: %w", err)
	}
	if err = zipFile.Close(); err != nil {
		return fmt.Errorf("failed to close zip file: %w", err)
	}

	filePath = zipFile.Name()
	status = 0
	return nil
}

func tryPeer(peersAddr []Addr, full bool) *Addr {
	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	type result struct {
		addr *Addr
	}
	resultCh := make(chan result, len(peersAddr))

	for _, addr := range peersAddr {
		addrCopy := addr
		go func() {
			udpAddr, err := net.ResolveUDPAddr("udp", addrCopy.String())
			if err != nil {
				fmt.Println("Failed to resolve address:", err)
				resultCh <- result{addr: nil}
				return
			}

			conn, err := quic.DialAddr(ctx, udpAddr.String(), &tls.Config{InsecureSkipVerify: true}, nil)
			fmt.Println("Trying to connect to:", addrCopy.String())
			if err != nil {
				fmt.Println("Failed to connect:", err)
				resultCh <- result{addr: nil}
				return
			}

			response := handleRequest("hel chalaoshi-next-client", conn)

			var prefix string
			if full {
				prefix = "hel chalaoshi-next-server-full"
			} else {
				prefix = "hel chalaoshi-next-server"
			}

			if strings.HasPrefix(response, prefix) {
				resultCh <- result{addr: &addrCopy}
			}
			resultCh <- result{addr: nil}
		}()
	}

	for i := 0; i < len(peersAddr); i++ {
		res := <-resultCh
		if res.addr != nil {
			cancel()
			return res.addr
		}
	}

	return nil
}
