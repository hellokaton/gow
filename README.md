# gow
🙄 gow!!! the micro web framework.

```go
g := gow.Me()
g.Get("/", func(ctx gow.Context) {
    log.Warn("Hello World")
    ctx.Text("Hello World \n")
}).Before(func(ctx gow.Context) {
    log.Debug("global before")
}).After(func(ctx gow.Context) {
    log.Debug("global after")
}).NotFound(func(ctx gow.Context) {
    ctx.Status(404)
    ctx.Text("not found this url")
})
g.Listen(":10777")
```