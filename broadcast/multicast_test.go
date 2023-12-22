package broadcast_test

import (
	"net"
	"pan/broadcast"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMulticast ...
func TestMulticast(t *testing.T) {

	t.Run("write and read success", func(t *testing.T) {

		buf := []byte("content buffer")

		serve, err := broadcast.NewMulticast()
		if err != nil {
			t.Fatal(err)
		}
		defer serve.Close()

		wg := sync.WaitGroup{}
		wg.Add(1)

		var rmsg []byte
		var hostport []byte
		go (func() {
			msg, addr, err := serve.Read(-1)
			wg.Done()
			if err != nil {
				t.Fatal(err)
			}
			rmsg = msg
			hostport = addr
		})()

		serve.Write(buf)
		wg.Wait()

		ips := make([]string, 0)
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			t.Fatal(err)
		}
		for _, addr := range addrs {
			naddr := addr.(*net.IPNet)
			ips = append(ips, naddr.IP.String())
		}

		ip, _, err := net.SplitHostPort(string(hostport))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, buf, rmsg, "Buffer should be same")
		assert.Contains(t, ips, ip, "IP should be local")
	})

}
