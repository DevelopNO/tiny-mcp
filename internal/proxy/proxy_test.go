package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
	"github.com/noamohana/mini-mcp/internal/auth"
	"github.com/noamohana/mini-mcp/internal/router"
	"github.com/rs/zerolog"
)

func TestProxy_EchoIntegration(t *testing.T) {
	t.Parallel()
	backend, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()
	go func() {
		for {
			conn, err := backend.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				io.Copy(c, c)
			}(conn)
		}
	}()
	pol := router.PolicyMap{"teamA": backend.Addr().String()}
	file := t.TempDir() + "/pol.yaml"
	f := struct{ Policies router.PolicyMap `yaml:"policies"` }{pol}
	b, _ := yaml.Marshal(f)
	os.WriteFile(file, b, 0o600)
	r, err := router.NewRouter(file)
	if err != nil {
		t.Fatal(err)
	}
	s := auth.NewSigner("testkey")
	p := &Proxy{
		Listen:   "127.0.0.1:0",
		Signer:   s,
		Router:   r,
		Logger:   zerolog.Nop(),
	}
	ln, err := net.Listen("tcp", p.Listen)
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() { for { c, _ := ln.Accept(); go p.handle(c) } }()
	client, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	fmt.Fprintln(client, "teamA")
	scan := bufio.NewScanner(client)
	scan.Scan() // nonce
	scan.Scan() // token
	msg := "hello\n"
	client.Write([]byte(msg))
	buf := make([]byte, len(msg))
	client.Read(buf)
	if string(buf) != msg {
		t.Fatalf("echo failed: got %q", buf)
	}
}
