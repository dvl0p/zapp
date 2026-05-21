package main

import (
	"errors"
	"bytes"
	"log"
	"os"
	"strconv"
	"net"
	"encoding/hex"
)

type Protocol int

const (
	TCP Protocol = iota
	TCP6
	UDP
	UDP6
)

type socketStatus int

const (
	SocketStatusListen socketStatus = 10
)

type Inode int

type Socket struct {
	Proto Protocol
	Status socketStatus
	Inode Inode
	Host net.IP
	Port uint16
}

func newSocket(proto Protocol, status socketStatus, 
		inode Inode, host net.IP, port uint16) Socket {
	
	return Socket{
		Proto: proto,
		Status: status,
		Inode: inode,
		Host: host,
		Port: port,
	}
}

func discoverSockets() ([]Socket, error) {

	protos := []Protocol{TCP, UDP, TCP6, UDP6}
	var sockets []Socket

	for _, proto := range protos {
		newSockets, err := scanProcFile(proto)
		if err != nil {
			return nil, err
		}
		sockets = append(sockets, 
			filterSockets(newSockets, SocketStatusListen)...)
	}

	return sockets, nil
}

func filterSockets(sockets []Socket, status socketStatus) []Socket {

	filteredSockets := []Socket{}
	for i := range sockets {
		if sockets[i].Status == status {
			filteredSockets = append(filteredSockets, sockets[i])
		}
	}

	return filteredSockets
}

func scanProcFile(p Protocol) ([]Socket, error) {

	procFileContents, err := os.ReadFile(getProcPath(p))
	if err != nil {
		return nil, err
	}

	return parseProcTable(procFileContents, p), nil
}

func parseProcTable(procFileContents []byte, p Protocol) []Socket {

	sockets := []Socket{}

	lines := bytes.Split(procFileContents, []byte("\n"))

	for i := 1; i < len(lines)-1 ; i++ {
		fields :=  bytes.Fields(lines[i])

		status, err := strconv.ParseInt(string(fields[3]), 16, 0)
		if err != nil {
			log.Println("Could not convert status for line: ", lines[i])
			continue
		}

		inode, err := strconv.Atoi(string(fields[9]))
		if err != nil {
			log.Println("Could not convert inode for line: ", lines[i])
			continue
		}

		hostHex, portHex, found := bytes.Cut(fields[1], []byte(":"))
		if !found {
			log.Println("Could not parse host and port from: ", 
				string(lines[i]))
			continue
		}
		
		host, err := hexToIP(hostHex, p)
		if err != nil {
			log.Println("Could not convert from hex to IP: ", lines[i])
		}

		port, err := strconv.ParseUint(string(portHex), 16, 16)
		if err != nil {
			log.Println("Could not convert from string to Inode: ", portHex)
		}

		socket := newSocket(p, socketStatus(status), Inode(inode), 
			host, uint16(port))
		sockets = append(sockets, socket)
	}

	return sockets
}

func hexToIP (hexIP []byte, p Protocol) (net.IP, error) {
	switch p {
	case TCP, UDP:
		return parserIPv4(hexIP)
	case TCP6, UDP6:
		return parserIPv6(hexIP)
	default:
		return nil, errors.New("protocol not valid")
	}
}

func parserIPv4(hexIP []byte) (net.IP, error) {
	b, err := hex.DecodeString(string(hexIP))
	if err != nil {
		return nil, err
	}

	return net.IPv4(b[3], b[2], b[1], b[0]), nil
}

func parserIPv6(hexIP []byte) (net.IP, error) {

    b, err := hex.DecodeString(string(hexIP))
	if err != nil {
		return nil, err
	}

    for i := 0; i < 16; i += 4 {
        b[i], b[i+3] = b[i+3], b[i]
        b[i+1], b[i+2] = b[i+2], b[i+1]
    }

    return net.IP(b), nil
}

func getProcPath(p Protocol) string {
	procNet := "/proc/net/"
	switch p {
	case TCP:
		return procNet + "tcp"
	case TCP6:
		return procNet + "tcp6"
	case UDP:
		return procNet + "udp"
	case UDP6:
		return procNet + "udp6"
	default:
		return ""
	}
}
