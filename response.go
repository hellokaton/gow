package gow

import (
	"net/http"
	"strconv"
	"time"
	"encoding/json"
	"os"
	"path"
	"encoding/xml"
)

type Response struct {
	http.ResponseWriter
	ctx         *context
	Status      int
	IsInterrupt bool
	IsCommit    bool
	Body        []byte
	Headers     map[string]string
	gow         *Gow
}

func NewResponse(w http.ResponseWriter) *Response {
	res := new(Response)
	res.ResponseWriter = w
	res.Status = http.StatusOK
	res.IsCommit = false
	res.IsInterrupt = false
	res.Headers = make(map[string]string)
	res.Headers["Content-Type"] = "text/html;charset=UTF-8"
	return res
}

// ContentType sets content-type string.
func (res *Response) ContentType(contentType string) {
	res.Headers["Content-Type"] = contentType
}

func (res *Response) Text(str string) {
	res.ContentType("application/json;charset=UTF-8")
	res.Body = []byte(str)
}

// Json set json response with data and proper header.
func (res *Response) Json(data interface{}) {
	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	res.ContentType("application/json;charset=UTF-8")
	res.Body = bytes
}

// Xml render
func (res *Response) Xml(data interface{}) {
	bytes, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	res.ContentType("application/xml")
	res.Body = bytes
}

// Cookie gets cookie value by given key when give only string.
// Cookie sets cookie value by given key, value and expire time string.
func (res *Response) Cookie(key string, value ...string) string {
	t := time.Now()
	expire, _ := strconv.Atoi(value[1])
	t = t.Add(time.Duration(expire) * time.Second)
	cookie := &http.Cookie{
		Name:    key,
		Value:   value[0],
		Path:    "/",
		MaxAge:  expire,
		Expires: t,
	}
	http.SetCookie(res, cookie)
	return ""
}

func (res *Response) Interrupt() {
	if res.IsInterrupt {
		return
	}
	if !res.IsCommit {
		res.Commit()
	}
	res.IsInterrupt = true
}

func (res *Response) Commit() {
	if res.IsCommit {
		return
	}
	for name, value := range res.Headers {
		res.Header().Set(name, value)
	}
	if res.IsInterrupt {
		res.WriteHeader(res.Status)
	}
	res.Write(res.Body)
	res.IsCommit = true
}

func (res *Response) Throw(status int) {
	res.Status = status
	res.Commit()
}

// Download sends file download response by file path.
func (res *Response) Download(file string) {
	f, e := os.Stat(file)
	if e != nil {
		res.Status = http.StatusNotFound
		return
	}
	if f.IsDir() {
		res.Status = http.StatusForbidden
		return
	}
	output := res.Header()
	output.Set("Content-Type", "application/octet-stream")
	output.Set("Content-Disposition", "attachment; filename="+path.Base(file))
	output.Set("Content-Transfer-Encoding", "binary")
	output.Set("Expires", "0")
	output.Set("Cache-Control", "must-revalidate")
	output.Set("Pragma", "public")
	http.ServeFile(res, res.ctx.request, file)
}

// Redirect does redirection response to url string and status int optional.
func (res *Response) Redirect(url string, status ...int) {
	res.Headers["Location"] = url
	if len(status) > 0 {
		res.Status = status[0]
		return
	}
	res.Status = http.StatusFound
}
