package peer_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"pan/core"
	"pan/peer"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTLSConf(isClient bool) (*tls.Config, *x509.Certificate) {

	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值

	template := &x509.Certificate{
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject: pkix.Name{ // 证书的主题信息
			Country:            []string{"CN"},         // 证书所属的国家
			Organization:       []string{"company"},    // 证书存放的公司名称
			OrganizationalUnit: []string{"department"}, // 证书所属的部门名称
			Province:           []string{"ChengDu"},    // 证书签发机构所在省
			CommonName:         "localhost",            // 证书域名
			Locality:           []string{"ChengDu"},    // 证书签发机构所在市
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, // 典型用法是指定叶子证书中的公钥的使用目的。它包括一系列的OID，每一个都指定一种用途。例如{id pkix 31}表示用于服务器端的TLS/SSL连接；{id pkix 34}表示密钥可以用于保护电子邮件。
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,                      // 指定了这份证书包含的公钥可以执行的密码操作，例如只能用于签名，但不能用来加密
		IsCA:                  true,                                                                       // 指示证书是不是ca证书
		BasicConstraintsValid: true,                                                                       // 指示证书是不是ca证书
	}

	// 生成公私钥对
	caPrivkey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	// pubKey, err := x509.MarshalPKIXPublicKey(&caPrivkey.PublicKey)
	// if err != nil {
	// 	panic(err)
	// }

	// 生成自签证书(template=parent)
	rootCertDer, err := x509.CreateCertificate(rand.Reader, template, template, &caPrivkey.PublicKey, caPrivkey) //DER 格式
	if err != nil {
		panic(err)
	}

	// 将私钥编码为pkcs8格式
	caPrivBytes, err := x509.MarshalPKCS8PrivateKey(caPrivkey)
	if err != nil {
		panic(err)
	}

	key := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: caPrivBytes})
	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootCertDer})

	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	x509Cert, err := core.ParseCertWithPem(cert)
	if err != nil {
		panic(err)
	}

	if isClient == true {
		tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}

		return tlsConf, x509Cert
	}
	tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}

	return tlsConf, x509Cert

}

// TestQuic node test cases
func TestQuic(t *testing.T) {

	t.Run("Serve and Dial success", func(t *testing.T) {
		// test connect
		serveTlsConf, serveCert := newTLSConf(false)
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
		if err != nil {
			t.Fatal(err)
		}

		var node peer.Node
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		serve, err := peer.ServeQUICNode(addr, serveTlsConf)
		if err != nil {
			t.Fatal(err)
		}
		defer serve.Close()
		go func() {
			n, err := serve.Accept(ctx)
			wg.Done()
			if err != nil {
				t.Fatal(err)
			}
			node = n
		}()

		clientTlsConf, clientCert := newTLSConf(true)
		timeOutCtx, timeOutCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer timeOutCancel()
		dialNode, err := peer.DialQUICNode(addr, clientTlsConf, timeOutCtx)
		if err != nil {
			t.Fatal(err)
		}
		defer dialNode.Close()

		wg.Wait()

		nodeCert := node.Certificate()
		dialNodeCert := dialNode.Certificate()

		assert.NotNil(t, node, "Node should not be nil")

		assert.Equal(t, serveCert.PublicKey, dialNodeCert.PublicKey, "Serve public key should be same")
		assert.Equal(t, serveCert.Signature, dialNodeCert.Signature, "Serve signature should be same")
		assert.Equal(t, clientCert.PublicKey, nodeCert.PublicKey, "Client public key should be same")
		assert.Equal(t, clientCert.Signature, nodeCert.Signature, "Client signature should be same")

		// test send and recv by stream
		wg = sync.WaitGroup{}
		wg.Add(1)
		var stream peer.NodeStream

		go func() {
			ctx := context.Background()
			stream, err = node.AcceptNodeStream(ctx)
			wg.Done()
			if err != nil {
				t.Fatal(err)
			}
		}()

		dialStream, err := dialNode.OpenNodeStream()
		if err != nil {
			t.Fatal(err)
		}

		content := "Request Content"
		io.WriteString(dialStream, content)

		err = dialStream.Close()
		if err != nil {
			t.Fatal(err)
		}

		wg.Wait()

		assert.NotNil(t, stream, "Accept Stream should not be nil")
		// assert.Equal(t, (*stream).StreamID(), (*dialStream).StreamID(), "Stream ID should be same")

		buf, err := io.ReadAll(stream)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, content, string(buf), "Stream request content should be same")

		resContent := "Response Content"
		io.WriteString(stream, resContent)

		err = stream.Close()
		if err != nil {
			t.Fatal(err)
		}

		resBuf, err := io.ReadAll(dialStream)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, resContent, string(resBuf), "Stream response content should be same")

	})

	t.Run("MarshalQUICAddr and UnmarshalQUICAddr", func(t *testing.T) {
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
		if err != nil {
			t.Fatal(err)
		}

		payload := peer.MarshalQUICAddr(addr)
		quicAddr, err := peer.UnmarshalQUICAddr(payload)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, addr.Network(), quicAddr.Network(), "Network should be udp")
		assert.Equal(t, addr.String(), quicAddr.String(), "Addr should be same")

	})

	t.Run("NodeDialer", func(t *testing.T) {
		serveTlsConf, _ := newTLSConf(false)
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
		if err != nil {
			t.Fatal(err)
		}

		var node peer.Node
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		serve, err := peer.ServeQUICNode(addr, serveTlsConf)
		if err != nil {
			t.Fatal(err)
		}
		defer serve.Close()
		go func() {
			n, err := serve.Accept(ctx)
			wg.Done()
			if err != nil {
				t.Fatal(err)
			}
			node = n
		}()

		clientTlsConf, _ := newTLSConf(true)
		timeOutCtx, timeOutCancel := context.WithTimeout(context.Background(), time.Second*5)
		dialer := peer.NewNodeDialer(clientTlsConf, timeOutCtx)
		defer timeOutCancel()

		quicAddr := peer.MarshalQUICAddr(addr)
		dialNode, err := dialer.Connect(quicAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer dialNode.Close()

		wg.Wait()
		assert.NotNil(t, node, "Node should not be nil")

	})

	t.Run("NodeStream CloseRead", func(t *testing.T) {
		serveTlsConf, _ := newTLSConf(false)
		addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
		if err != nil {
			t.Fatal(err)
		}

		var node peer.Node
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		serve, err := peer.ServeQUICNode(addr, serveTlsConf)
		if err != nil {
			t.Fatal(err)
		}
		defer serve.Close()
		go func() {
			n, err := serve.Accept(ctx)
			wg.Done()
			if err != nil {
				t.Fatal(err)
			}
			node = n
		}()

		clientTlsConf, _ := newTLSConf(true)
		timeOutCtx, timeOutCancel := context.WithTimeout(context.Background(), time.Second*5)
		dialer := peer.NewNodeDialer(clientTlsConf, timeOutCtx)
		defer timeOutCancel()

		quicAddr := peer.MarshalQUICAddr(addr)
		dialNode, err := dialer.Connect(quicAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer dialNode.Close()

		wg.Wait()

		// test send and recv by stream
		wg = sync.WaitGroup{}
		wg.Add(1)
		var stream peer.NodeStream

		go func() {
			defer wg.Done()
			ctx := context.Background()
			stream, err = node.AcceptNodeStream(ctx)
			if err != nil {
				t.Fatal(err)
			}

		}()

		dialStream, err := dialNode.OpenNodeStream()
		if err != nil {
			t.Fatal(err)
		}

		io.WriteString(dialStream, "reader header")
		wg.Wait()

		err = stream.CloseRead()
		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(dialStream, bytes.NewReader([]byte("reader header")))
		if err != nil {
			t.Fatal(err)
		}

		// wg.Wait()

		assert.NotNil(t, stream, "Accept Stream should not be nil")

	})

}
