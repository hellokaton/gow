package main

import (
	"github.com/biezhi/gow"
	"html/template"
)

func IndexPage(ctx *gow.Context) string {
	return "Hello Gow!"
}

func Hello(ctx gow.Context) {
	ctx.Text("Hello World\n")
}

func main() {
	g := gow.Me()
	
	g.Get("/", func(ctx gow.Context) {
		ctx.Text("Hello World \n")
	}).Get("/hello", func(ctx gow.Context) {
		fmap := template.FuncMap{
			"Name": "gow",
		}
		ctx.Render("hello", fmap)
	}).Get("/users", func(ctx gow.Context) {
		fmap := template.FuncMap{
			"u1": "jack1",
			"u2": "jack2",
			"u3": "jack3",
		}
		ctx.Json(fmap)
	}).Before(func(ctx gow.Context) {
		g.Logger().Debug("global before")
	}).After(func(ctx gow.Context) {
		g.Logger().Debug("global after")
	}).NotFound(func(ctx gow.Context) {
		g.Logger().Debug("404: %s", ctx.Uri())
		ctx.Status(404)
		ctx.Text("not found this url")
	})
	g.Use(func(ctx gow.Context) {
		g.Logger().Debug("中间件执行...")
	})
	g.Listen()
}
