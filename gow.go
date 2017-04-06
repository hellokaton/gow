package gow

import (
	"net/http"
	"github.com/biezhi/agon/log"
	"sync"
	"crypto/tls"
	"net"
)

type Gow struct {
	Router
	view         *View
	DisableHTTP2 bool
	Dev          bool
	Server       *http.Server
	TLSServer    *http.Server
	Listener     net.Listener
	TLSListener  net.Listener
	pool         sync.Pool
}

func Me() *Gow {
	gow := Gow{}
	gow.inter = make(map[string]Handler)
	statics = append(statics, "public")
	statics = append(statics, "static")
	gow.ApplyStatic()
	return &gow
}

func (g *Gow) enableSSL(address string) error {
	s := g.TLSServer
	s.Addr = address
	if !g.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	return g.StartServer(g.TLSServer)
}

// StartServer starts a custom http server.
func (e *Gow) StartServer(s *http.Server) (err error) {
	// Setup
	s.Handler = e
	
	if s.TLSConfig == nil {
		if e.Listener == nil {
			e.Listener, err = newListener(s.Addr)
			if err != nil {
				return err
			}
		}
		log.Info("⇛ http server started on %s\n", e.Listener.Addr())
		return s.Serve(e.Listener)
	}
	if e.TLSListener == nil {
		l, err := newListener(s.Addr)
		if err != nil {
			return err
		}
		e.TLSListener = tls.NewListener(l, s.TLSConfig)
	}
	log.Info("⇛ http server started on %s\n", e.TLSListener.Addr())
	return s.Serve(e.TLSListener)
}

func (g *Gow) Listen(addr string) {
	log.Info("Server Listen %s", addr)
	log.Error("Server Start Error", http.ListenAndServe(addr, g))
}

func (g *Gow) View() *View {
	return g.view
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func newListener(address string) (*tcpKeepAliveListener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}
