package core

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const BufSize = 1024

type SecureSocket struct {
	Cipher     *Cipher
	ListAddr   *net.TCPAddr
	RemoteAddr *net.TCPAddr
}

func (SecureSocket *SecureSocket) DecodeRead(conn io.ReadWriter, bs []byte) (n int, err error) {
	n, err = conn.Read(bs)
	if err != nil {
		return
	}
	SecureSocket.Cipher.decode(bs[:n])
	return
}

func (SecureSocket *SecureSocket) EncodeWrite(conn io.ReadWriter, bs []byte) (int, error) {
	SecureSocket.Cipher.encode(bs)
	return conn.Write(bs)
}

func (SecureSocket *SecureSocket) EncodeCopy(dst io.ReadWriter, src io.ReadWriter) error {
	buf := make([]byte, BufSize)
	for {
		readCount, errRead := src.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := SecureSocket.EncodeWrite(dst, buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

func (SecureSocket *SecureSocket) DecodeCopy(dst io.ReadWriter, src io.ReadWriter) error {
	buf := make([]byte, BufSize)
	for {
		readCount, err := SecureSocket.DecodeRead(src, buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, err := dst.Write(buf[0:readCount])
			if err != nil {
				return err
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}

// 和远程的socket建立连接，他们之间的数据传输会加密
func (secureSocket *SecureSocket) DialRemote() (*net.TCPConn, error) {
	remoteConn, err := net.DialTCP("tcp", nil, secureSocket.RemoteAddr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("连接到远程服务器 %s 失败:%s", secureSocket.RemoteAddr, err))
	}
	return remoteConn, nil
}
