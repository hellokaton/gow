package gow

import (
	"net/http"
	"strconv"
	"net/url"
	"strings"
	"path"
)

type (
	Context interface {
		Base() string
		Header(key string, value ...string) string
		Host() string
		Ip() string
		Url() string
		Uri() string
		Ext() string
		Method() string
		UserAgent() string
		Referer() string
		IsSSL() bool
		IsAjax() bool
		Request() *http.Request
		Response() *Response
		Text(text string)
		Param(key string) string
		QueryParam(key string) string
		Input() map[string]string
		String(key string) string
		Strings(key string) []string
		Flash(key string, v ...interface{}) interface{}
		
		SetBody(body []byte)
		Status(code int)
		ContentType(contentType string)
		Render(tpl string, model map[string]interface{})
		Tpl(tpl string, data map[string]interface{}) string
	}
	
	context struct {
		request   *http.Request
		response  *Response
		path      string
		pnames    []string
		pvalues   []string
		query     url.Values
		gow       *Gow
		flashData map[string]interface{}
	}
)

// Flash sets values to this context or gets by key string.
// The flash items are alive in this context only.
func (ctx *context) Flash(key string, v ...interface{}) interface{} {
	if len(v) > 0 {
		return ctx.flashData[key]
	}
	ctx.flashData[key] = v[0]
	return nil
}

// Input returns all input data map.
func (ctx *context) Input() map[string]string {
	data := make(map[string]string)
	for key, v := range ctx.request.Form {
		data[key] = v[0]
	}
	return data
}

// Strings returns string slice of given key.
func (ctx *context) Strings(key string) []string {
	return ctx.request.Form[key]
}

// String returns input value of given key.
func (ctx *context) String(key string) string {
	return ctx.request.FormValue(key)
}

// StringOr returns input value of given key instead of def string if empty.
func (ctx *context) StringOr(key string, def string) string {
	value := ctx.String(key)
	if value == "" {
		return def
	}
	return value
}

// Int returns input value of given key.
func (ctx *context) Int(key string) int {
	str := ctx.String(key)
	i, _ := strconv.Atoi(str)
	return i
}

// IntOr returns input value of given key instead of def int if empty.
func (ctx *context) IntOr(key string, def int) int {
	i := ctx.Int(key)
	if i == 0 {
		return def
	}
	return i
}

// Float returns input value of given key.
func (ctx *context) Float(key string) float64 {
	str := ctx.String(key)
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

// FloatOr returns input value of given key instead of def float if empty.
func (ctx *context) FloatOr(key string, def float64) float64 {
	f := ctx.Float(key)
	if f == 0.0 {
		return def
	}
	return f
}

// Bool returns input value of given key.
func (ctx *context) Bool(key string) bool {
	str := ctx.String(key)
	b, _ := strconv.ParseBool(str)
	return b
}

// ContentType sets content-type string.
func (ctx *context) ContentType(contentType string) {
	ctx.response.ContentType(contentType)
}

func (ctx *context) Tpl(tpl string, data map[string]interface{}) string {
	//b, e := ctx.gow.TplEngine.Render(tpl+".html", data)
	//if e != nil {
	//	panic(e)
	//}
	//return string(b)
	return ""
}

func (ctx *context) Render(tpl string, model map[string]interface{}) {
	b, e := _gow.TplEngine.Render(tpl, model)
	if e != nil {
		panic(e)
	}
	ctx.ContentType("text/html; charset=utf-8")
	ctx.response.Body = b
}

// Gow returns *Gow instance in this context.
func (ctx *context) Gow() *Gow {
	return ctx.gow
}

// Gow returns *Gow instance in this context.
func (ctx *context) Text(text string) {
	ctx.response.Text(text)
}

func (ctx *context) Request() *http.Request {
	return ctx.request
}

func (ctx *context) Response() *Response {
	return ctx.response
}

func (ctx *context) IsAjax() bool {
	return ctx.request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

func (ctx *context) IsSSL() bool {
	return ctx.request.TLS != nil
}

func (ctx *context) Ip() string {
	return strings.Split(ctx.request.RemoteAddr, ":")[0]
}

func (ctx *context) Referer() string {
	return ctx.request.Referer()
}

func (ctx *context) UserAgent() string {
	return ctx.request.UserAgent()
}

func (ctx *context) Method() string {
	return ctx.request.Method
}

func (ctx *context) Url() string {
	return ctx.request.URL.String()
}

func (ctx *context) Uri() string {
	return ctx.request.RequestURI
}

func (ctx *context) Ext() string {
	return path.Ext(ctx.request.URL.Path)
}

func (ctx *context) Host() string {
	return ctx.request.Host
}

func (ctx *context) Header(key string, value ...string) string {
	if len(value) == 0 {
		return ctx.request.Header.Get(key)
	}
	if len(value) == 1 {
		ctx.response.Headers[key] = value[0]
	}
	return ""
}

func (ctx *context) Base() string {
	baseUrl := "://" + ctx.Host() + "/"
	if ctx.IsSSL() {
		baseUrl = "https" + baseUrl
	} else {
		baseUrl = "http" + baseUrl
	}
	return baseUrl
}

func (ctx *context) Param(key string) string {
	return ""
}

func (ctx *context) QueryParam(key string) string {
	return ctx.request.Form.Get(key)
}

func (ctx *context) SetBody(body []byte) {
	ctx.response.Body = body
}

func (ctx *context) Status(code int) {
	ctx.response.Status = code
}
