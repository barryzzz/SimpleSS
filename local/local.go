package local

import (
	"log"
	"net"

	"com.lsl/ss/core"
)

type LsLocal struct {
	*core.SecureSocket
}

func NewLsLocal(password *core.Password, listAddr, remoteAdd *net.TCPAddr) *LsLocal {
	return &LsLocal{
		SecureSocket: &core.SecureSocket{
			Cipher:     core.NewCipher(password),
			ListAddr:   listAddr,
			RemoteAddr: remoteAdd,
		},
	}
}

func (local *LsLocal) Listen(didListen func(listenAddr net.Addr)) error {
	listener, err := net.ListenTCP("tcp", local.ListAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	if didListen != nil {
		didListen(listener.Addr())
	}
	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		userConn.SetLinger(0)
		go local.handleConn(userConn)
	}
}

func (local *LsLocal) handleConn(conn *net.TCPConn) {
	defer conn.Close()
	log.Println("收到请求-->", conn.RemoteAddr().String())
	proxyServer, err := local.DialRemote()
	if err != nil {
		log.Println(err)
		return
	}
	defer proxyServer.Close()
	proxyServer.SetLinger(0)

	go func() {
		err := local.DecodeCopy(conn, proxyServer)
		if err != nil {
			conn.Close()
			proxyServer.Close()
		}
	}()
	local.EncodeCopy(proxyServer, conn)
}
