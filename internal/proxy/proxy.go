package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/noamohana/mini-mcp/internal/auth"
	"github.com/noamohana/mini-mcp/internal/router"
	"github.com/rs/zerolog"
)

type Proxy struct {
	Listen   string
	Signer   *auth.Signer
	Router   *router.Router
	Logger   zerolog.Logger
}

func (p *Proxy) Serve() error {
	ln, err := net.Listen("tcp", p.Listen)
	if err != nil {
		return err
	}
	p.Logger.Info().Str("listen", p.Listen).Msg("proxy started")
	for {
		conn, err := ln.Accept()
		if err != nil {
			p.Logger.Error().Err(err).Msg("accept failed")
			continue
		}
		go p.handle(conn)
	}
}

func (p *Proxy) handle(conn net.Conn) {
	if conn == nil {
		return
	}
	defer conn.Close()
	peer := conn.RemoteAddr().String()
	r := bufio.NewReader(conn)
	teamLine, err := r.ReadString('\n')
	if err != nil {
		p.Logger.Error().Err(err).Str("peer", peer).Msg("read team failed")
		return
	}
	team := strings.TrimSpace(teamLine)
	backend, ok := p.Router.Lookup(team)
	if !ok {
		p.Logger.Warn().Str("team", team).Msg("unknown team")
		return
	}
	nonce := time.Now().UnixNano()
	token, err := p.Signer.Sign(peer, team)
	if err != nil {
		p.Logger.Error().Err(err).Msg("sign failed")
		return
	}
	bconn, err := net.Dial("tcp", backend)
	if err != nil {
		p.Logger.Error().Err(err).Str("backend", backend).Msg("dial backend failed")
		return
	}
	defer bconn.Close()
	// handshake: send nonce and token
	fmt.Fprintf(conn, "%d\n%s\n", nonce, token)
	// pipe streams
	done := make(chan struct{}, 2)
	go func() { io.Copy(bconn, r); done <- struct{}{} }()
	go func() { io.Copy(conn, bconn); done <- struct{}{} }()
	<-done
	<-done
	p.Logger.Info().Str("peer", peer).Str("team", team).Str("backend", backend).Msg("session closed")
}
