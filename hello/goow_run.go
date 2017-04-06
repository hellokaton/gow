package main

import (
	"github.com/biezhi/gow"
	"fmt"
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
	}).Before(func(ctx gow.Context) {
		fmt.Println("global before")
	}).After(func(ctx gow.Context) {
		fmt.Println("global after")
	}).NotFound(func(ctx gow.Context) {
		fmt.Println("404...")
		ctx.Status(404)
		ctx.Text("not found this url")
	})
	g.Listen()
}
