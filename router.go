package gow

import (
	nurl "net/url"
	"regexp"
	"path"
	"strings"
	"github.com/biezhi/agon/log"
	"net/http"
	"fmt"
	"time"
	"sync"
	"github.com/biezhi/agon/utils"
	"runtime/debug"
)

const (
	RequestMethodANY     = RequestMethod("ANY")
	RequestMethodGet     = RequestMethod("GET")
	RequestMethodPost    = RequestMethod("POST")
	RequestMethodPut     = RequestMethod("PUT")
	RequestMethodDelete  = RequestMethod("DELETE")
	RequestMethodPatch   = RequestMethod("PATCH")
	RequestMethodHead    = RequestMethod("HEAD")
	RequestMethodOptions = RequestMethod("OPTIONS")
	
	NOT_FOUND   = "not_found"
	STATIC_RES  = "static_resource"
	BEFORE_EXEC = "before_exec"
	AFTER_EXEC  = "after_exec"
)

var statics []string

type (
	RequestMethod string
	
	Route struct {
		pattern *regexp.Regexp
		method  RequestMethod
		params  []string
		handler []Handler
	}
	
	Router struct {
		routes []*Route
		inter  map[string]Handler
		middle []Handler
		Lock   sync.RWMutex
	}
	
	Handler func(Context)
)

// Find does find matched rule and parse route url, returns route params and matched handlers.
func (rt *Router) Find(url string, method string) (params map[string]string, fn []Handler) {
	sfx := path.Ext(url)
	url = strings.Replace(url, sfx, "", -1)
	// fix path end slash
	url = nurl.QueryEscape(url)
	
	//if !strings.HasSuffix(url, "%2F") && sfx == "" {
	//	url += "%2F"
	//	log.Debug("aaaa")
	//}
	
	url = strings.Replace(url, "%2F", "/", -1)
	
	for _, r := range rt.routes {
		if r.pattern.MatchString(url) && string(r.method) == method {
			p := r.pattern.FindStringSubmatch(url)
			if len(p) != len(r.params)+1 {
				continue
			}
			params = make(map[string]string)
			for i, n := range r.params {
				params[n] = p[i+1]
			}
			fn = r.handler
			return
		}
	}
	return nil, nil
}

func (r *Router) Get(pattern string, fn ...Handler) *Router {
	return r.Add(RequestMethodGet, pattern, fn...)
}

func (r *Router) Post(pattern string, fn ...Handler) *Router {
	return r.Add(RequestMethodPost, pattern, fn...)
}

func (r *Router) Put(pattern string, fn ...Handler) *Router {
	return r.Add(RequestMethodPut, pattern, fn...)
}

func (r *Router) Delete(pattern string, fn ...Handler) *Router {
	return r.Add(RequestMethodDelete, pattern, fn...)
}

func (r *Router) NotFound(h Handler) *Router {
	r.inter[NOT_FOUND] = h
	return r
}

func (r *Router) Static(v ...string) *Router {
	for _, s := range v {
		statics = append(statics, s)
	}
	utils.RemoveDuplicatesAndEmpty(statics)
	r.ApplyStatic()
	return r
}

func (r *Router) ApplyStatic() {
	r.inter[STATIC_RES] = func(ctx Context) {
		url := strings.TrimPrefix(ctx.Uri(), "/")
		for _, s := range statics {
			if strings.HasPrefix(url, s) {
				http.ServeFile(ctx.Response(), ctx.Request(), url)
				ctx.Response().Finish()
			}
		}
		
	}
}

func (r *Router) Before(h Handler) *Router {
	r.inter[BEFORE_EXEC] = h
	return r
}

func (r *Router) After(h Handler) *Router {
	r.inter[AFTER_EXEC] = h
	return r
}

func (r *Router) Add(method RequestMethod, pattern string, fn ...Handler) *Router {
	regex, params := ParsePattern(pattern)
	r.routes = append(r.routes, &Route{
		method:  method,
		params:  params,
		pattern: regex,
		handler: fn,
	})
	log.Info("Add Route %s\t%s", string(method), pattern)
	return r
}

func ParsePattern(pattern string) (regex *regexp.Regexp, params []string) {
	params = make([]string, 0)
	segments := strings.Split(nurl.QueryEscape(pattern), "%2F")
	for i, v := range segments {
		if strings.HasPrefix(v, "%3A") {
			segments[i] = `([\w-%]+)`
			params = append(params, strings.TrimPrefix(v, "%3A"))
		}
	}
	regex, _ = regexp.Compile("^" + strings.Join(segments, "/") + "$")
	return
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Acquire lock
	r.Lock.RLock()
	defer r.Lock.RUnlock()
	r.handler(w, req)
}

func (r *Router) Use(m ...Handler) *Router {
	r.middle = append(r.middle, m...)
	return r
}

func NewContext(w http.ResponseWriter, req *http.Request) Context {
	req.ParseForm()
	return &context{
		request:   req,
		response:  NewResponse(w),
		flashData: make(map[string]interface{}),
	}
}

func (r *Router) handler(w http.ResponseWriter, req *http.Request) {
	context := NewContext(w, req)
	
	defer func() {
		e := recover()
		if e == nil {
			context = nil
			return
		}
		context.SetBody([]byte(fmt.Sprint(e)))
		context.Status(http.StatusServiceUnavailable)
		println(string(context.Response().Body))
		debug.PrintStack()
		
		if _, ok := r.inter[AFTER_EXEC]; ok {
			r.inter[AFTER_EXEC](context)
		}
		if !context.Response().IsFinish {
			context.Response().Finish()
		}
		context = nil
	}()
	
	if _, ok := r.inter[STATIC_RES]; ok {
		r.inter[STATIC_RES](context)
		if context.Response().IsFinish {
			return
		}
	}
	
	if _, ok := r.inter[BEFORE_EXEC]; ok {
		r.inter[BEFORE_EXEC](context)
		if context.Response().IsFinish {
			return
		}
	}
	
	if len(r.middle) > 0 {
		for _, h := range r.middle {
			h(context)
			if context.Response().IsFinish {
				break
			}
		}
	}
	
	if context.Response().IsCommit {
		return
	}
	
	var (
		params map[string]string
		fn     []Handler
		url    = req.URL.Path
	)
	params, fn = r.Find(url, req.Method)
	if params != nil && fn != nil {
		start := time.Now()
		//context.routeParams = params
		for _, f := range fn {
			f(context)
			if context.Response().IsFinish {
				break
			}
		}
		if !context.Response().IsCommit {
			context.Response().Commit()
		}
		log.Info("%s\t%s\t%s", string(req.Method), req.URL.Path, time.Since(start))
	} else {
		context.Response().Status = 404
		if _, ok := r.inter[NOT_FOUND]; ok {
			r.inter[NOT_FOUND](context)
			if !context.Response().IsFinish {
				context.Response().Finish()
			}
		} else {
			log.Warn("Not Found %s", req.URL.Path)
			context.Response().Throw(404)
		}
	}
	if _, ok := r.inter[AFTER_EXEC]; ok {
		r.inter[AFTER_EXEC](context)
	}
	context = nil
}
