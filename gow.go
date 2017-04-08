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
	DisableHTTP2 bool
	Dev          bool
	Server       *http.Server
	TLSServer    *http.Server
	Listener     net.Listener
	TLSListener  net.Listener
	Pool         sync.Pool
	Host         string
	Port         int
	TplEngine    *TemplateEngine
}

var logger = log.NewLog()
var _gow = Gow{}

func Me() *Gow {
	_gow.Init()
	return &_gow
}

func (g *Gow) Init() {
	logger.Prefix("[gow] -- ")
	_gow.inter = make(map[string]Handler)
	_gow.TplEngine = NewTemplateEngine()
	_gow.Host = "0.0.0.0"
	_gow.Port = 10077
	statics = append(statics, "public")
	statics = append(statics, "static")
	statics = append(statics, "favicon.ico")
	_gow.Config()
	_gow.ApplyStatic()
}

func (g *Gow) Logger() *log.Logger {
	return logger
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
			logger.Prefix(prefix + " ")
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
func (g *Gow) StartServer(s *http.Server) (err error) {
	// Setup
	s.Handler = g
	
	if s.TLSConfig == nil {
		if g.Listener == nil {
			g.Listener, err = newListener(s.Addr)
			if err != nil {
				return err
			}
		}
		logger.Info("⇛ http server started on %s\n", g.Listener.Addr())
		return s.Serve(g.Listener)
	}
	if g.TLSListener == nil {
		l, err := newListener(s.Addr)
		if err != nil {
			return err
		}
		g.TLSListener = tls.NewListener(l, s.TLSConfig)
	}
	logger.Info("⇛ http server started on %s\n", g.TLSListener.Addr())
	return s.Serve(g.TLSListener)
}

func (g *Gow) Listen(addr ... string) {
	if len(addr) > 0 {
		logger.Info("Server Listen %s", addr)
		logger.Error("Server Start Error", http.ListenAndServe(addr[0], g))
	} else {
		logger.Info("Server Listen %s:%d", g.Host, g.Port)
		address := g.Host + ":" + fmt.Sprint(g.Port)
		logger.Error("Server Start Error", http.ListenAndServe(address, g))
	}
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
