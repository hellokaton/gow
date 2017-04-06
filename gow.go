package gow

import (
	"net/http"
	"github.com/biezhi/agon/log"
	"github.com/biezhi/agon/conf/ini"
	"sync"
	"crypto/tls"
	"net"
	"fmt"
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
	Pool         sync.Pool
	Host         string
	Port         int
}

var logger = log.NewLog()

func Me() *Gow {
	gow := Gow{}
	logger.Prefix("[gow] -- ")
	gow.inter = make(map[string]Handler)
	gow.Host = "0.0.0.0"
	gow.Port = 10077
	statics = append(statics, "public")
	statics = append(statics, "static")
	gow.Config()
	gow.ApplyStatic()
	return &gow
}

func (g *Gow) Config() *Gow {
	file, err := ini.LoadFile("app.ini")
	if err != nil {
		logger.Error("load [app.ini] error %s", err.Error())
	} else {
		host, ok := file.Get("app", "host")
		if ok {
			g.Host = host
		}
		port, ok := file.GetInt("app", "port")
		if ok {
			g.Port = port
		}
		prefix, ok := file.Get("logger", "prefix")
		if ok {
			logger.Prefix(prefix)
		}
		dev, ok := file.GetBool("app", "dev")
		if ok {
			g.Dev = dev
		}
	}
	return g
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
		logger.Info("⇛ http server started on %s\n", e.Listener.Addr())
		return s.Serve(e.Listener)
	}
	if e.TLSListener == nil {
		l, err := newListener(s.Addr)
		if err != nil {
			return err
		}
		e.TLSListener = tls.NewListener(l, s.TLSConfig)
	}
	logger.Info("⇛ http server started on %s\n", e.TLSListener.Addr())
	return s.Serve(e.TLSListener)
}

func (g *Gow) Listen(addr ... string) {
	if len(addr) > 0 {
		logger.Info("Server Listen %s", addr)
		logger.Error("Server Start Error", http.ListenAndServe(addr[0], g))
	} else{
		logger.Info("Server Listen %s:%d", g.Host, g.Port)
		address := g.Host + ":" + fmt.Sprint(g.Port)
		logger.Error("Server Start Error", http.ListenAndServe(address, g))
	}
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
