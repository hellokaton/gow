package main

import (
	"net/http"
	"strings"
	"fmt"
)

//返回一个Router实例
func New() *Router {
	return new(Router)
}

//路由结构体，包含一个记录方法、路径的map
type Router struct {
	Route map[string]map[string]http.HandlerFunc
}

//实现Handler接口，匹配方法以及路径
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h, ok := r.Route[req.Method][req.URL.String()]; ok {
		h(w, req)
	}
}

// 根据方法、路径将方法注册到路由
func (r *Router) HandleFunc(method, path string, f http.HandlerFunc) {
	method = strings.ToUpper(method)
	if r.Route == nil {
		r.Route = make(map[string]map[string]http.HandlerFunc)
	}
	if r.Route[method] == nil {
		r.Route[method] = make(map[string]http.HandlerFunc)
	}
	r.Route[method][path] = f
}

func main() {
	r := New()
	r.HandleFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello Get!")
	})
	r.HandleFunc("POST", "/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello POST!")
	})
	http.ListenAndServe(":10777", r)
}