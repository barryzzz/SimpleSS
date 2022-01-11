package server

import (
	"encoding/binary"
	"log"
	"net"

	"com.lsl/ss/core"
)

type LsServer struct {
	*core.SecureSocket
}

func NewServer(password string, listAddr string) (*LsServer, error) {
	pwd, err := core.ParsePassword(password)
	if err != nil {
		return nil, err
	}
	serverListen, err := net.ResolveTCPAddr("tcp", listAddr)
	if err != nil {
		return nil, err
	}

	return &LsServer{
		SecureSocket: &core.SecureSocket{
			Cipher:   core.NewCipher(pwd),
			ListAddr: serverListen,
		},
	}, nil
}

func (LsServer *LsServer) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", LsServer.ListAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	if didListen != nil {
		didListen(listener.Addr())
	}
	for {
		localConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		localConn.SetLinger(0)
		go LsServer.handleConn(localConn)
	}

}

func (LsServer *LsServer) handleConn(localConn *net.TCPConn) {
	defer localConn.Close()
	buf := make([]byte, 256)
	_, err := LsServer.DecodeRead(localConn, buf)
	if err != nil || buf[0] != 0x05 {
		return
	}
	LsServer.EncodeWrite(localConn, []byte{0x05, 0x00})
	n, err := LsServer.DecodeRead(localConn, buf)
	if err != nil || n < 7 {
		return
	}
	if buf[1] != 0x01 {
		return
	}

	var dIP []byte
	switch buf[3] {
	case 0x01:
		dIP = buf[4 : 4+net.IPv4len]
	case 0x03:
		ipAddr, err := net.ResolveIPAddr("ip", string(buf[5:n-2]))
		if err != nil {
			log.Println(err)
			return
		}
		dIP = ipAddr.IP
	case 0x04:
		dIP = buf[4 : 4+net.IPv6len]
	default:
		return
	}
	dPort := buf[n-2:]
	dstAddr := &net.TCPAddr{
		IP:   dIP,
		Port: int(binary.BigEndian.Uint16(dPort)),
	}
	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	log.Println("监控:", dstAddr.IP)
	if err != nil {
		log.Println(err)
		return
	} else {
		defer dstServer.Close()
		dstServer.SetLinger(0)
		LsServer.EncodeWrite(localConn, []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	}

	go func() {
		err := LsServer.DecodeCopy(dstServer, localConn)
		if err != nil {
			log.Println(err)
			localConn.Close()
			dstServer.Close()
		}
	}()

	LsServer.EncodeCopy(localConn, dstServer)

}
