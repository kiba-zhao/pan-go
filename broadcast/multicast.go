package broadcast

import (
	"errors"

	"net"
)

type multicastSt struct {
	serve      *net.UDPConn
	dialer     *net.UDPConn
	bufferSize int
}

// Read ...
func (b *multicastSt) Read(size int) (payload []byte, srcAddr []byte, err error) {
	var bufferSize int
	if size <= 0 {
		bufferSize = b.bufferSize
	} else {
		bufferSize = size
	}
	buff := make([]byte, bufferSize)
	byteLen, addr, err := b.serve.ReadFromUDP(buff)
	if err != nil {
		return
	}
	if byteLen == len(buff) {
		payload = buff
	} else {
		payload = buff[:byteLen]
	}

	srcAddr = []byte(addr.String())
	return
}

// Write ...
func (b *multicastSt) Write(payload []byte) (err error) {

	size, err := b.dialer.Write(payload)

	if err != nil {
		return
	}
	if size < len(payload) {
		err = &NetWriteError{writenSize: size}
	}

	return
}

// Close ...
func (b *multicastSt) Close() error {
	serr := b.serve.Close()
	derr := b.dialer.Close()
	return errors.Join(serr, derr)
}

type newMulticastConfig struct {
	addr       *net.UDPAddr
	ifi        *net.Interface
	bufferSize int
}

// defaultMulticastConf ...
func defaultMulticastConfig() (*newMulticastConfig, error) {

	addr, err := net.ResolveUDPAddr("udp", "224.0.0.120:9100")
	if err != nil {
		return nil, err
	}

	cfg := new(newMulticastConfig)
	cfg.addr = addr
	cfg.bufferSize = 64 * 1024 * 1024
	return cfg, nil
}

type NewMulticastWithFn func(cfg *newMulticastConfig)

// WithNewMulticastWithAddr ...
func NewMulticastWithAddr(addr *net.UDPAddr) NewMulticastWithFn {
	return func(cfg *newMulticastConfig) {
		cfg.addr = addr
	}
}

func NewMulticast(withFns ...NewMulticastWithFn) (Net, error) {

	cfg, err := defaultMulticastConfig()
	if err != nil {
		return nil, err
	}

	for _, withFn := range withFns {
		withFn(cfg)
	}

	serve, err := net.ListenMulticastUDP(cfg.addr.Network(), cfg.ifi, cfg.addr)
	if err != nil {
		return nil, err
	}
	dialer, err := net.DialUDP(cfg.addr.Network(), nil, cfg.addr)
	if err != nil {
		return nil, err
	}

	n := new(multicastSt)
	n.dialer = dialer
	n.serve = serve
	n.bufferSize = cfg.bufferSize
	return n, nil
}
